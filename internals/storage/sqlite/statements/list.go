package statements

import (
	"context"
	"database/sql"
)

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
