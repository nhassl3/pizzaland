package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/domain/models"
	"github.com/nhassl3/pizzaland/internals/lib/logger/sl"
	"github.com/nhassl3/pizzaland/internals/lib/marshall"
	"github.com/nhassl3/pizzaland/internals/storage"
	"github.com/nhassl3/pizzaland/internals/storage/sqlite/statements"
)

const (
	opSave               = "sqlite.Save"
	opSaveCategory       = "sqlite.SaveCategory"
	opGet                = "sqlite.Get"
	opGetCategory        = "sqlite.GetCategory"
	opList               = "sqlite.List"
	opRemove             = "sqlite.Remove"
	opRemoveCategory     = "sqlite.RemoveCategory"
	opUpdate             = "sqlite.Update"
	opUpdateCategoryById = "sqlite.UpdateCategoryById"
	opUpdateCategoryName = "sqlite.UpdateCategoryName"
)

type Storage struct {
	st *statements.Statement
}

func NewStorage(timeout time.Duration, path string) (*Storage, error) {
	db, err := sql.Open(
		"sqlite3",
		path+fmt.Sprintf("?_timeout=%d&_journal=WAL&_sync=NORMAL&_cache=shared&_fk=true&_txlock=immediate", timeout.Milliseconds()))
	if err != nil {
		return nil, err
	}

	// Ограничиваем количество открытых соединений
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0) // infinity

	st := statements.NewStatement(db)

	return &Storage{st}, nil
}

func (s *Storage) Save(ctx context.Context, pizza *pizzalndv1.PizzaProperties) (pizzaId uint64, err error) {
	typeDoughEnums := pizza.GetTypeDough()
	typeDoughInts := make([]int32, len(typeDoughEnums))
	for i, dough := range typeDoughEnums {
		typeDoughInts[i] = int32(dough) // Преобразование enum в базовый int32
	}

	pizzaId, err = s.st.Save(
		ctx,
		pizza.GetCategoryId(),
		pizza.GetName(),
		pizza.GetDescription().GetValue(),
		typeDoughInts,
		pizza.GetPrice(),
		pizza.GetDiameter(),
	)
	var sqliteErr sqlite3.Error
	if err != nil {
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, sl.ErrUpLevel(opSave, storage.ErrPizzaExists)
		}
		return 0, sl.ErrUpLevel(opSave, err)
	}

	return
}

func (s *Storage) SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (categoryId uint32, err error) {
	res, err := s.st.SaveCategory(
		ctx,
		category.GetName(),
		category.GetDescription().GetValue(),
	)
	var sqliteErr sqlite3.Error
	if err != nil {
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, sl.ErrUpLevel(opSaveCategory, storage.ErrPizzaExists)
		}
		return 0, sl.ErrUpLevel(opSaveCategory, err)
	}

	if id, err := res.RowsAffected(); err == nil {
		categoryId = uint32(id)
	} else {
		return 0, sl.ErrUpLevel(opSaveCategory, err)
	}

	return
}

func (s *Storage) Get(ctx context.Context, ident any) (pizza *pizzalndv1.PizzaProperties, err error) {
	pizzaObj, err := s.st.Get(ctx, ident)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opGet, storage.ErrPizzaNotFound)
		}
		return nil, sl.ErrUpLevel(opGet, err)
	}

	destPizza := &pizzalndv1.PizzaProperties{}
	pizza, err = marshall.MarshalModels(pizzaObj, destPizza)
	if err != nil {
		return nil, sl.ErrUpLevel(opGet, err)
	}

	return
}

func (s *Storage) GetCategory(ctx context.Context, ident any) (category *pizzalndv1.CategoryProperties, err error) {
	var categoryObj models.Category

	if err := s.st.GetCategory(ctx, ident, &categoryObj); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opGetCategory, storage.ErrCategoryNotFound)
		}
		return nil, sl.ErrUpLevel(opGetCategory, err)
	}

	destCategory := &pizzalndv1.CategoryProperties{}

	if category, err = marshall.MarshalModels(&categoryObj, destCategory); err != nil {
		return nil, sl.ErrUpLevel(opGetCategory, err)
	}

	return
}

func (s *Storage) List(ctx context.Context, ident any, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	pizzas, err := s.st.List(ctx, ident, offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opList, storage.ErrPizzaNotFound)
		}
		return nil, sl.ErrUpLevel(opList, err)
	}

	if len(pizzas) == 0 {
		return nil, sl.ErrUpLevel(opList, storage.ErrListPizzaOutOfRange)
	}

	for i := range pizzas {
		destPizza := &pizzalndv1.PizzaProperties{}
		destPizza, err = marshall.MarshalModels(&pizzas[i], destPizza)
		if err != nil {
			return nil, sl.ErrUpLevel(opGet, err)
		}
		pizza = append(pizza, destPizza)
	}

	return
}

func (s *Storage) Remove(ctx context.Context, ident any) (success bool, err error) {
	success, err = s.st.Remove(ctx, ident)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opRemove, storage.ErrPizzaNotFound)
		}
		return false, sl.ErrUpLevel(opRemove, err)
	}

	return
}

func (s *Storage) RemoveCategory(ctx context.Context, ident any) (success bool, err error) {
	success, err = s.st.RemoveCategory(ctx, ident)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opRemoveCategory, storage.ErrCategoryNotFound)
		}
		return false, sl.ErrUpLevel(opRemoveCategory, err)
	}

	return
}

func (s *Storage) Update(
	ctx context.Context,
	ident any,
	categoryId uint32,
	name string,
	description string,
	typeDough []pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (bool, error) {
	var sqliteErr sqlite3.Error

	typeDoughInt32 := make([]int32, len(typeDough))
	for i, v := range typeDough {
		typeDoughInt32[i] = int32(v)
	}

	if err := s.st.Update(ctx, ident, categoryId, name, description, typeDoughInt32, price, diameter); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opUpdate, storage.ErrPizzaNotFound)
		} else if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return false, sl.ErrUpLevel(opUpdate, storage.ErrPizzaExists)
		}
		return false, sl.ErrUpLevel(opUpdate, err)
	}

	return true, nil
}

func (s *Storage) UpdateCategoryById(ctx context.Context, id uint32, name string, description string) (success bool, err error) {
	res, err := s.st.UpdateCategory(ctx, id, name, description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opUpdateCategoryById, storage.ErrPizzaNotFound)
		}
		return false, sl.ErrUpLevel(opUpdateCategoryById, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return false, sl.ErrUpLevel(opUpdateCategoryById, err)
	}

	return true, nil
}

func (s *Storage) UpdateCategoryByName(ctx context.Context, name string, description string) (success bool, err error) {
	res, err := s.st.UpdateCategory(ctx, 0, name, description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opUpdateCategoryName, storage.ErrPizzaNotFound)
		}
		return false, sl.ErrUpLevel(opUpdateCategoryName, err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return false, sl.ErrUpLevel(opUpdateCategoryName, err)
	}

	return true, nil
}
