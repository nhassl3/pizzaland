package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/domain/models"
	"github.com/nhassl3/pizzaland/internals/lib/logger/sl"
	"github.com/nhassl3/pizzaland/internals/lib/marshall"
	"github.com/nhassl3/pizzaland/internals/storage"
)

const (
	opSave               = "sqlite.Save"
	opSaveCategory       = "sqlite.SaveCategory"
	opGetById            = "sqlite.GetById"
	opGetByName          = "sqlite.GetByName"
	opGetCategoryById    = "sqlite.GetCategoryById"
	opGetCategoryByName  = "sqlite.GetCategoryByName"
	opList               = "sqlite.List"
	opListCategoryById   = "sqlite.ListCategoryById"
	opListCategoryByName = "sqlite.ListCategoryByName"
	opRemoveById         = "sqlite.RemoveById"
	opRemoveByName       = "sqlite.RemoveByName"
	opRemoveCategoryById = "sqlite.RemoveCategoryById"
	opRemoveCategoryName = "sqlite.RemoveCategoryName"
	opUpdateById         = "sqlite.UpdateById"
	opUpdateByName       = "sqlite.UpdateByName"
	opUpdateCategoryById = "sqlite.UpdateCategoryById"
	opUpdateCategoryName = "sqlite.UpdateCategoryName"
)

type Storage struct {
	st *Statement
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	st := NewStatement(db)

	return &Storage{st}, nil
}

func (s *Storage) Save(ctx context.Context, pizza *pizzalndv1.PizzaProperties) (pizzaId uint64, err error) {
	res, err := s.st.Save(
		ctx,
		pizza.GetCategoryId(),
		pizza.GetName(),
		pizza.GetDescription().GetValue(),
		*pizza.GetTypeDough()[0].Enum(),
		pizza.GetPrice(),
		pizza.GetDiameter(),
	)
	var sqliteErr sqlite3.Error
	if err != nil {
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, sl.ErrUpLevel(opSave, storage.ErrPizzaExists.Error())
		}
		return 0, sl.ErrUpLevel(opSave, err.Error())
	}

	if id, err := res.LastInsertId(); err == nil {
		pizzaId = uint64(id)
	} else {
		return 0, sl.ErrUpLevel(opSave, err.Error())
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
			return 0, sl.ErrUpLevel(opSaveCategory, storage.ErrPizzaExists.Error())
		}
		return 0, sl.ErrUpLevel(opSaveCategory, err.Error())
	}

	if id, err := res.LastInsertId(); err == nil {
		categoryId = uint32(id)
	} else {
		return 0, sl.ErrUpLevel(opSaveCategory, err.Error())
	}

	return
}

func (s *Storage) GetById(ctx context.Context, id uint64) (pizza *pizzalndv1.PizzaProperties, err error) {
	var pizzaObj models.Pizza

	res, err := s.st.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opGetById, storage.ErrPizzaNotFound.Error())
		}
		return nil, sl.ErrUpLevel(opGetById, err.Error())
	}

	if err := res.Scan(
		&pizzaObj.PizzaId,
		&pizzaObj.CategoryId,
		&pizzaObj.Name,
		&pizzaObj.Description,
		&pizzaObj.TypeDough,
		&pizzaObj.Price,
		&pizzaObj.Diameter,
	); err != nil {
		return nil, sl.ErrUpLevel(opGetById, err.Error())
	}

	destPizza := &pizzalndv1.PizzaProperties{}
	pizza, err = marshall.MarshalModels(&pizzaObj, destPizza)
	if err != nil {
		return nil, sl.ErrUpLevel(opGetById, err.Error())
	}

	return
}

func (s *Storage) GetByName(ctx context.Context, name string) (pizza *pizzalndv1.PizzaProperties, err error) {
	res, err := s.st.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opGetByName, storage.ErrPizzaNotFound.Error())
		}
		return nil, sl.ErrUpLevel(opGetByName, err.Error())
	}

	if err := res.Scan(pizza); err != nil {
		return nil, sl.ErrUpLevel(opGetByName, err.Error())
	}

	return
}

func (s *Storage) GetCategoryById(ctx context.Context, id uint32) (category *pizzalndv1.CategoryProperties, err error) {
	res, err := s.st.GetCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opGetCategoryById, storage.ErrPizzaNotFound.Error())
		}
		return nil, sl.ErrUpLevel(opGetCategoryById, err.Error())
	}

	if err := res.Scan(category); err != nil {
		return nil, sl.ErrUpLevel(opGetCategoryById, err.Error())
	}

	return
}

func (s *Storage) GetCategoryByName(ctx context.Context, name string) (category *pizzalndv1.CategoryProperties, err error) {
	res, err := s.st.GetCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opGetCategoryByName, storage.ErrPizzaNotFound.Error())
		}
		return nil, sl.ErrUpLevel(opGetCategoryByName, err.Error())
	}

	if err := res.Scan(category); err != nil {
		return nil, sl.ErrUpLevel(opGetCategoryByName, err.Error())
	}

	return
}

func (s *Storage) List(ctx context.Context, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	res, err := s.st.List(ctx, "", offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opList, storage.ErrPizzaNotFound.Error())
		}
		return nil, sl.ErrUpLevel(opList, err.Error())
	}

	if err := res.Scan(pizza); err != nil {
		return nil, sl.ErrUpLevel(opList, err.Error())
	}

	return
}

func (s *Storage) ListCategoryByName(ctx context.Context, categoryName string, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	res, err := s.st.List(ctx, categoryName, offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opListCategoryByName, storage.ErrPizzaNotFound.Error())
		}
		return nil, sl.ErrUpLevel(opListCategoryByName, err.Error())
	}

	if err := res.Scan(pizza); err != nil {
		return nil, sl.ErrUpLevel(opListCategoryByName, err.Error())
	}

	return
}

func (s *Storage) ListCategoryById(ctx context.Context, categoryId, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	res, err := s.st.ListByCategoryId(ctx, categoryId, offset, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sl.ErrUpLevel(opListCategoryById, storage.ErrPizzaNotFound.Error())
		}
		return nil, sl.ErrUpLevel(opListCategoryById, err.Error())
	}

	if err := res.Scan(pizza); err != nil {
		return nil, sl.ErrUpLevel(opListCategoryById, err.Error())
	}

	return
}

func (s *Storage) RemoveById(ctx context.Context, id uint64) (success bool, err error) {
	res, err := s.st.RemoveById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opRemoveById, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opRemoveById, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opRemoveById, err.Error())
	}

	slog.Info(opRemoveById, slog.Int64("LastInsertId", lastInsertId))

	return int64(id) == lastInsertId, nil
}

func (s *Storage) RemoveByName(ctx context.Context, name string) (success bool, err error) {
	res, err := s.st.RemoveByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opRemoveByName, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opRemoveByName, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opRemoveByName, err.Error())
	}

	slog.Info(opRemoveByName, slog.Int64("LastInsertId", lastInsertId))

	return true, nil
}

func (s *Storage) RemoveCategoryById(ctx context.Context, id uint32) (success bool, err error) {
	res, err := s.st.RemoveCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opRemoveCategoryById, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opRemoveCategoryById, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opRemoveCategoryById, err.Error())
	}

	slog.Info(opRemoveCategoryById, slog.Int64("LastInsertId", lastInsertId))

	return int64(id) == lastInsertId, nil
}

func (s *Storage) RemoveCategoryByName(ctx context.Context, name string) (success bool, err error) {
	res, err := s.st.RemoveCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opRemoveCategoryName, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opRemoveCategoryName, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opRemoveCategoryName, err.Error())
	}

	slog.Info(opRemoveCategoryName, slog.Int64("LastInsertId", lastInsertId))

	return true, nil
}

func (s *Storage) UpdateById(
	ctx context.Context,
	id uint64,
	categoryId uint32,
	name string,
	description string,
	typeDough *pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (success bool, err error) {
	res, err := s.st.Update(ctx, id, categoryId, name, description, typeDough, price, diameter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opUpdateById, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opUpdateById, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opUpdateById, err.Error())
	}

	slog.Info(opUpdateById, slog.Int64("LastInsertId", lastInsertId))

	return int64(id) == lastInsertId, nil
}

func (s *Storage) UpdateByName(
	ctx context.Context,
	categoryId uint32,
	name string,
	description string,
	typeDough *pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (success bool, err error) {
	res, err := s.st.Update(ctx, 0, categoryId, name, description, typeDough, price, diameter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opUpdateByName, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opUpdateByName, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opUpdateByName, err.Error())
	}

	slog.Info(opUpdateByName, slog.Int64("LastInsertId", lastInsertId))

	return true, nil
}

func (s *Storage) UpdateCategoryById(ctx context.Context, id uint32, name string, description string) (success bool, err error) {
	res, err := s.st.UpdateCategory(ctx, id, name, description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opUpdateCategoryById, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opUpdateCategoryById, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opUpdateCategoryById, err.Error())
	}

	slog.Info(opUpdateCategoryById, slog.Int64("LastInsertId", lastInsertId))

	return int64(id) == lastInsertId, nil
}

func (s *Storage) UpdateCategoryByName(ctx context.Context, name string, description string) (success bool, err error) {
	res, err := s.st.UpdateCategory(ctx, 0, name, description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, sl.ErrUpLevel(opUpdateCategoryName, storage.ErrPizzaNotFound.Error())
		}
		return false, sl.ErrUpLevel(opUpdateCategoryName, err.Error())
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return false, sl.ErrUpLevel(opUpdateCategoryName, err.Error())
	}

	slog.Info(opUpdateCategoryName, slog.Int64("LastInsertId", lastInsertId))

	return true, nil
}
