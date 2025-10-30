package sqlite

import (
	"context"
	"database/sql"
	"strings"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
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
	typeDough pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(
		ctx,
		"INSERT INTO pizza (category_id, name, description, type_dough, price, diameter) VALUES (?, ?, ?, ?, ?, ?);",
	)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, categoryId, name, description, typeDough, price, diameter)
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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, name, description)
}

func (s *Statement) GetById(
	ctx context.Context,
	id uint64,
) (*sql.Row, error) {
	stmt, err := s.db.PrepareContext(ctx,
		"SELECT * FROM pizza WHERE id=?;")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.QueryRowContext(ctx, id), nil
}

// GetByName returns result with select request
func (s *Statement) GetByName(
	ctx context.Context,
	name string,
) (*sql.Row, error) {
	stmt, err := s.db.PrepareContext(ctx,
		"SELECT * FROM pizza WHERE name=?;")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.QueryRowContext(ctx, name), nil
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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.QueryRowContext(ctx, name), nil
}

// List returns the pizza list with by default category or if category is set with category
func (s *Statement) List(ctx context.Context, categoryName string, offset, limit uint32) (*sql.Rows, error) {
	var query string
	var args []any

	if categoryName == "" {
		query = "SELECT * FROM pizza OFFSET ? LIMIT ?;"
		args = []any{offset, limit}
	} else {
		query = "SELECT * FROM pizza WHERE category_id=(SELECT id FROM categories WHERE name=?) LIMIT ? OFFSET ?;"
		args = []any{categoryName, offset, limit}
	}

	stmt, err := s.db.PrepareContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.QueryContext(ctx, categoryId, offset, limit)
}

// RemoveById removes pizza record from the system by id of the pizza
func (s *Statement) RemoveById(ctx context.Context, id uint64) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM pizza WEHRE id=?")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, id)
}

// RemoveByName deletes pizza record from the system by name of the pizza (recommend)
func (s *Statement) RemoveByName(ctx context.Context, name string) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM pizza WEHRE name=?")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, name)
}

// RemoveCategoryById removes category from the system and removes everyone pizza in this category list by CASCADE
func (s *Statement) RemoveCategoryById(ctx context.Context, id uint32) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM categories WEHRE id=?")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, id)
}

// RemoveCategoryByName removes category from the system and removes everyone pizza in this category list by CASCADE
func (s *Statement) RemoveCategoryByName(ctx context.Context, name string) (sql.Result, error) {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM categories WEHRE name=?")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, name)
}

// Update updating pizza table on the system
// Also checks user specifies more than one parameter
func (s *Statement) Update(ctx context.Context,
	id uint64,
	categoryId uint32,
	name string,
	description string,
	typeDough *pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (sql.Result, error) {
	var query string
	args := make([]interface{}, 0, 6)
	updates := make([]string, 0, 6)

	if categoryId != 0 {
		updates = append(updates, "category_id=?")
		args = append(args, categoryId)
	}

	if name != "" {
		updates = append(updates, "name=?")
		args = append(args, name)
	}

	if description != "" {
		updates = append(updates, "description=?")
		args = append(args, description)
	}

	if *typeDough != pizzalndv1.TypeDough_UNKNOWN {
		updates = append(updates, "type_dough=?")
		args = append(args, int(typeDough.Number()))
	}

	if price != 0 {
		updates = append(updates, "price=?")
		args = append(args, price)
	}

	if diameter != 0 {
		updates = append(updates, "diameter=?")
		args = append(args, diameter)
	}

	if len(updates) == 0 {
		return nil, storage.ErrNothingToChangePizza
	}

	if id != 0 {
		query = "UPDATE pizza SET " + strings.Join(updates, ", ") + " WHERE id=?;"
		args = append(args, id)
	} else {
		query = "UPDATE pizza SET " + strings.Join(updates, ", ") + " WHERE name=?;"
		args = append(args, name)
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, args...)
}

// UpdateCategory update a category by name or ID specified in the parameters and check if the field value is equal to zero or nil.
// Also, an update occurs when the user specifies more than one of the parameters
func (s *Statement) UpdateCategory(ctx context.Context, id uint32, name, descriptions string) (sql.Result, error) {
	var query string
	updates := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)

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
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			panic(err)
		}
	}(stmt)

	return stmt.ExecContext(ctx, args...)
}
