package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/nhassl3/pizzaland/internals/storage"
)

type Statement struct {
	db *sql.DB
}

func NewStatement(db *sql.DB) *Statement {
	return &Statement{
		db: db,
	}
}

func (s *Statement) Save(
	ctx context.Context,
	categoryId uint32,
	name, description string,
	doughPizzaTypes []int32,
	price float32,
	diameter uint32,
) (uint64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertPizza, err := tx.PrepareContext(ctx,
		"INSERT INTO pizza (category_id, name, description, price, diameter) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer insertPizza.Close()

	res, err := insertPizza.ExecContext(ctx, categoryId, name, description, price, diameter)
	if err != nil {
		return 0, fmt.Errorf("failed to insert pizza record: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	if err := s.saveTypeDough(ctx, tx, uint64(id), doughPizzaTypes); err != nil {
		return 0, fmt.Errorf("failed to save dough: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return uint64(id), nil
}

func (s *Statement) saveTypeDough(
	ctx context.Context,
	tx *sql.Tx,
	pizzaId uint64,
	doughPizzaTypes []int32,
) error {
	if len(doughPizzaTypes) == 0 {
		return nil
	}

	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, pizzaType := range doughPizzaTypes {
		if _, err := stmt.ExecContext(ctx, pizzaId, pizzaType); err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return nil
}

func (s *Statement) SaveCategory(
	ctx context.Context,
	name, description string,
) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx,
		"INSERT INTO categories (name, description) VALUES (?, ?);")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, name, description)
}

func (s *Statement) Get(
	ctx context.Context,
	ident any,
) (*sql.Row, error) {
	var query string

	switch ident.(type) {
	case string:
		query = "SELECT * FROM pizza WHERE name=?;"
	case uint64:
		query = "SELECT * FROM pizza WHERE id=?;"
	default:
		return nil, storage.ErrInvalidIdentifier
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, ident)

	return row, nil
}

func (s *Statement) GetTypeDough(
	ctx context.Context,
	ident any,
) (*sql.Rows, error) {
	var query string

	switch ident.(type) {
	case string:
		query = "SELECT dough_type_id FROM pizza_dough_types WHERE pizza_id=(SELECT id FROM pizza where name=?)"
	case uint64:
		query = "SELECT dough_type_id FROM pizza_dough_types WHERE pizza_id=?"
	default:
		return nil, storage.ErrInvalidIdentifier
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, ident)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// GetCategoryById returns
func (s *Statement) GetCategoryById(
	ctx context.Context,
	id uint32,
) (*sql.Row, error) {
	stmt, err := s.db.PrepareContext(ctx,
		"SELECT * FROM categories WHERE id=?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, id), nil
}

// GetCategoryByName returns category
func (s *Statement) GetCategoryByName(
	ctx context.Context,
	name string,
) (*sql.Row, error) {
	stmt, err := s.db.PrepareContext(ctx,
		"SELECT * FROM categories WHERE name=?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, name), nil
}

// List returns the pizza list with by default category or if category is set with category
func (s *Statement) List(ctx context.Context, categoryName string, offset, limit uint32) (*sql.Rows, error) {
	var query string
	var args []any

	if categoryName == "" {
		query = "SELECT * FROM pizza ORDER BY name DESC LIMIT ? OFFSET ?;"
		args = []any{offset, limit}
	} else {
		query = "SELECT * FROM pizza WHERE category_id=(SELECT id FROM categories WHERE name=?) ORDER BY name DESC LIMIT ? OFFSET ?;"
		args = []any{categoryName, offset, limit}
	}

	stmt, err := s.db.PrepareContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}

// ListByCategoryId returns the pizza list with by default category or if category is set with category
func (s *Statement) ListByCategoryId(ctx context.Context, categoryId, offset, limit uint32) (*sql.Rows, error) {
	stmt, err := s.db.PrepareContext(
		ctx,
		"SELECT * FROM pizza WHERE category_id=? OFFSET ? LIMIT ?;",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, categoryId, offset, limit)
}

// RemoveById removes pizza record from the system by id of the pizza
func (s *Statement) RemoveById(ctx context.Context, id uint64) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM pizza WEHRE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, id)
}

// RemoveByName deletes pizza record from the system by name of the pizza (recommend)
func (s *Statement) RemoveByName(ctx context.Context, name string) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM pizza WEHRE name=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, name)
}

// RemoveCategoryById removes category from the system and removes everyone pizza in this category list by CASCADE
func (s *Statement) RemoveCategoryById(ctx context.Context, id uint32) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM categories WEHRE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, id)
}

// RemoveCategoryByName removes category from the system and removes everyone pizza in this category list by CASCADE
func (s *Statement) RemoveCategoryByName(ctx context.Context, name string) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM categories WEHRE name=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, name)
}

// Update updates pizza table and dough types in transaction
func (s *Statement) Update(ctx context.Context,
	ident any,
	categoryId uint32,
	name string,
	description string,
	typeDough []int32,
	price float32,
	diameter uint32,
) error {
	// 1. Подготавливаем данные для обновления основной таблицы
	updates, args := s.preparePizzaUpdates(categoryId, name, description, price, diameter)
	if len(updates) == 0 {
		return storage.ErrNothingToChangePizza
	}

	// 2. Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 3. Обновляем основную таблицу
	if err := s.updatePizzaTable(ctx, tx, ident, updates, args); err != nil {
		return err
	}

	// 4. Обновляем типы теста
	if err := s.updateDoughTypes(ctx, tx, ident, typeDough); err != nil {
		return err
	}

	// 5. Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateCategory update a category by name or ID specified in the parameters and check if the field value is equal to zero or nil.
// Also, an update occurs when the user specifies more than one of the parameters
func (s *Statement) UpdateCategory(ctx context.Context, id uint32, name, descriptions string) (sql.Result, error) {
	var query string
	updates := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if name != "" {
		updates = append(updates, "name=?")
		args = append(args, name)
	}

	if descriptions != "" {
		updates = append(updates, "description=?")
		args = append(args, descriptions)
	}

	if len(updates) == 0 {
		return nil, storage.ErrNothingToChangeCategory
	}

	if id != 0 {
		query = "UPDATE categories SET " + strings.Join(updates, ", ") + " WHERE id=?;"
		args = append(args, id)
	} else {
		query = "UPDATE categories SET " + strings.Join(updates, ", ") + " WHERE name=?;"
		args = append(args, name)
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, args...)
}

// preparePizzaUpdates подготавливает данные для обновления пиццы
func (s *Statement) preparePizzaUpdates(categoryId uint32, name, description string, price float32, diameter uint32) ([]string, []interface{}) {
	updates := make([]string, 0, 5)
	args := make([]interface{}, 0, 5)

	addUpdate := func(condition bool, column string, value interface{}) {
		if condition {
			updates = append(updates, column+"=?")
			args = append(args, value)
		}
	}

	addUpdate(categoryId != 0, "category_id", categoryId)
	addUpdate(name != "", "name", name)
	addUpdate(description != "", "description", description)
	addUpdate(price != 0, "price", price)
	addUpdate(diameter != 0, "diameter", diameter)

	return updates, args
}

// updatePizzaTable обновляет основную таблицу пицц
func (s *Statement) updatePizzaTable(ctx context.Context, tx *sql.Tx, ident any, updates []string, args []any) error {
	var query string

	switch ident.(type) {
	case string:
		query = "UPDATE pizza SET " + strings.Join(updates, ", ") + " WHERE name=?"
	case uint64:
		query = "UPDATE pizza SET " + strings.Join(updates, ", ") + " WHERE id=?"
	default:
		return storage.ErrInvalidIdentifier
	}

	args = append(args, ident)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare pizza update: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, args...); err != nil {
		return fmt.Errorf("execute pizza update: %w", err)
	}

	return nil
}

// updateDoughTypes обновляет типы теста для пиццы
func (s *Statement) updateDoughTypes(ctx context.Context, tx *sql.Tx, ident any, typeDough []int32) error {
	// 1. Удаляем старые типы теста
	if err := s.deleteDoughTypes(ctx, tx, ident); err != nil {
		return err
	}

	// 2. Если нет новых типов - выходим
	if len(typeDough) == 0 {
		return nil
	}

	// 3. Добавляем новые типы теста
	return s.insertDoughTypes(ctx, tx, ident, typeDough)
}

// deleteDoughTypes удаляет все типы теста для пиццы
func (s *Statement) deleteDoughTypes(ctx context.Context, tx *sql.Tx, ident any) error {
	var query string

	switch ident.(type) {
	case string:
		query = "DELETE FROM pizza_dough_types WHERE pizza_id = (SELECT id FROM pizza WHERE name = ?)"
	case uint64:
		query = "DELETE FROM pizza_dough_types WHERE pizza_id = ?"
	default:
		return storage.ErrInvalidIdentifier
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare dough types delete: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, ident); err != nil {
		return fmt.Errorf("execute dough types delete: %w", err)
	}

	return nil
}

// insertDoughTypes добавляет новые типы теста
func (s *Statement) insertDoughTypes(ctx context.Context, tx *sql.Tx, ident any, typeDough []int32) error {
	if len(typeDough) == 0 {
		err := tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
		return nil
	}
	var query string

	switch ident.(type) {
	case string:
		query = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES ((SELECT id FROM pizza WHERE name = ?), ?)"
	case uint64:
		query = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES (?, ?)"
	default:
		return storage.ErrInvalidIdentifier
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare dough types insert: %w", err)
	}
	defer stmt.Close()

	for _, doughTypeID := range typeDough {
		if _, err := stmt.ExecContext(ctx, ident, doughTypeID); err != nil {
			return fmt.Errorf("insert dough type %d: %w", doughTypeID, err)
		}
	}

	return nil
}
