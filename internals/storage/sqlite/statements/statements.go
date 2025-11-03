package statements

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/domain/models"
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
) (id uint64, err error) {
	var insertPizzaQuery, insertTypesDoughQuery string
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertPizzaQuery = "INSERT INTO pizza (category_id, name, description, price, diameter) VALUES (?, ?, ?, ?, ?)"
	insertTypesDoughQuery = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES (?, ?)"

	id, err = savePizza(ctx, tx, insertPizzaQuery, categoryId, name, description, price, diameter)
	if err != nil {
		return 0, fmt.Errorf("failed to save pizza: %w", err)
	}

	if err = saveTypeDough(ctx, tx, insertTypesDoughQuery, id, doughPizzaTypes); err != nil {
		return 0, fmt.Errorf("failed to save type dough: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return
}

// SaveCategory saves category of pizza in the system
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
) (*models.Pizza, error) {
	var pizzaObj models.Pizza
	var doughTypes []int32
	var getPizzaQuery, getTypeDoughQuery string

	switch ident.(type) {
	case string:
		getPizzaQuery = "SELECT * FROM pizza WHERE name=?;"
		getTypeDoughQuery = "SELECT dough_type_id FROM pizza_dough_types WHERE pizza_id=(SELECT id FROM pizza where name=?)"
	case uint64:
		getPizzaQuery = "SELECT * FROM pizza WHERE id=?;"
		getTypeDoughQuery = "SELECT dough_type_id FROM pizza_dough_types WHERE pizza_id=?"
	default:
		return nil, storage.ErrInvalidIdentifier
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	if err := getPizza(ctx, tx, getPizzaQuery, ident, &pizzaObj); err != nil {
		return nil, err
	}

	if err := getTypeDough(ctx, tx, getTypeDoughQuery, ident, &doughTypes); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	for _, v := range doughTypes {
		pizzaObj.TypeDough = append(pizzaObj.TypeDough, pizzalndv1.TypeDough(v))
	}

	return &pizzaObj, nil
}

// GetCategory returns info about category by identifier
func (s *Statement) GetCategory(
	ctx context.Context,
	ident any,
	categoryObj *models.Category,
) error {
	var query string

	switch ident.(type) {
	case string:
		query = "SELECT * FROM categories WHERE id=(SELECT id FROM categories where name=?);"
	case uint32:
		query = "SELECT * FROM categories WHERE id=?;"
	}
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, ident)

	if err := row.Scan(&categoryObj.Id, &categoryObj.Name, &categoryObj.Description); err != nil {
		return err
	}

	return nil
}

// List returns the pizza list with by default category or if category is set with category
func (s *Statement) List(ctx context.Context, ident any, offset, limit uint32) ([]models.Pizza, error) {
	var query string
	var args []any
	pizzas := make([]models.Pizza, 0, limit)

	if reflect.ValueOf(ident).IsZero() {
		query = "SELECT * FROM pizza ORDER BY id LIMIT ? OFFSET ?;"
		args = []any{limit, offset}
	} else {
		switch ident.(type) {
		case string:
			query = "SELECT * FROM pizza WHERE category_id=(SELECT id FROM categories WHERE name=?) ORDER BY id LIMIT ? OFFSET ?;"
		case uint32:
			query = "SELECT * FROM pizza WHERE category_id=? ORDER BY id LIMIT ? OFFSET ?;"
		default:
			return nil, storage.ErrInvalidIdentifier
		}
		args = []any{ident, limit, offset}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	if err := getPizzasSlice(ctx, tx, query, &pizzas, args); err != nil {
		return nil, fmt.Errorf("failed to retrieve pizzas: %w", err)
	}

	getTypeDoughQuery := "SELECT dough_type_id FROM pizza_dough_types WHERE pizza_id=?"
	for i := 0; i < len(pizzas); i++ {
		doughType := make([]int32, 0)
		if err := getTypeDough(ctx, tx, getTypeDoughQuery, pizzas[i].PizzaId, &doughType); err != nil {
			return nil, err
		}
		for _, doughTypeId := range doughType {
			pizzas[i].TypeDough = append(pizzas[i].TypeDough, pizzalndv1.TypeDough(doughTypeId))
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return pizzas, nil
}

// Update updates pizza table and dough types in transaction
func (s *Statement) Update(
	ctx context.Context,
	ident any,
	categoryId uint32,
	name string,
	description string,
	typeDough []int32,
	price float32,
	diameter uint32,
) error {
	var updateQuery, deleteTypesQuery, insertTypesQuery string

	updates, args := prepareDataUpdates(map[string]interface{}{
		"category_id": categoryId,
		"name":        name,
		"description": description,
		"price":       price,
		"diameter":    diameter,
	})
	if len(updates) == 0 {
		return storage.ErrNothingToChangePizza
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	switch ident.(type) {
	case string:
		updateQuery = "UPDATE pizza SET " + strings.Join(updates, ", ") + " WHERE name=?"
		deleteTypesQuery = "DELETE FROM pizza_dough_types WHERE pizza_id=(SELECT id FROM pizza WHERE name=?)"
		insertTypesQuery = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES ((SELECT id FROM pizza WHERE name=?), ?)"
	case uint64:
		updateQuery = "UPDATE pizza SET " + strings.Join(updates, ", ") + " WHERE id=?"
		deleteTypesQuery = "DELETE FROM pizza_dough_types WHERE pizza_id=?"
		insertTypesQuery = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES (?, ?)"
	}
	args = append(args, ident)

	if err := updatePizzaTable(ctx, tx, updateQuery, args); err != nil {
		return err
	}

	if err := updateDoughTypes(ctx, tx, deleteTypesQuery, insertTypesQuery, ident, typeDough); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateCategory update a category by name or ID specified in the parameters
// and check if the field value is equal to zero or nil.
// Also, an update occurs when the user specifies more than one of the parameters
func (s *Statement) UpdateCategory(ctx context.Context, ident any, name, descriptions string) (sql.Result, error) {
	var query string

	updates, args := prepareDataUpdates(map[string]interface{}{
		"name":        name,
		"description": descriptions,
	})
	if len(updates) == 0 {
		return nil, storage.ErrNothingToChangeCategory
	}

	switch ident.(type) {
	case string:
		query = "UPDATE categories SET " + strings.Join(updates, ", ") + " WHERE name=?;"
	case uint32:
		query = "UPDATE categories SET " + strings.Join(updates, ", ") + " WHERE id=?;"
	}
	args = append(args, ident)

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.ExecContext(ctx, args...)
}

// Remove removes pizza record from the system by identifier of the pizza
func (s *Statement) Remove(ctx context.Context, ident any) (bool, error) {
	var query string
	switch ident.(type) {
	case string:
		query = "DELETE FROM pizza WHERE name = ?"
	case uint64:
		query = "DELETE FROM pizza WHERE id = ?"
	default:
		return false, storage.ErrInvalidIdentifier
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, ident); err != nil {
		return false, err
	}

	return true, nil
}

// RemoveCategory removes category from the system and removes everyone pizza in this category list by CASCADE
func (s *Statement) RemoveCategory(ctx context.Context, ident any) (bool, error) {
	var query string

	switch ident.(type) {
	case string:
		query = "DELETE FROM categories WHERE name = ?"
	case uint32:
		query = "DELETE FROM categories WHERE id = ?"
	default:
		return false, storage.ErrInvalidIdentifier
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, ident); err != nil {
		return false, err
	}

	return true, nil
}
