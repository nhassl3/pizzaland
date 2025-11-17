package statements

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

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
	name, description, imagePath string,
	doughPizzaTypes []int32,
	rating uint32,
	price float32,
	sizes []int32,
) (id uint64, err error) {
	var insertPizzaQuery, insertTypesDoughQuery, insertSizesQuery string
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertPizzaQuery = "INSERT INTO pizza (category_id, name, description, price, rating, image_path) VALUES (?, ?, ?, ?, ?, ?)"
	insertTypesDoughQuery = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES (?, ?)"
	insertSizesQuery = "INSERT INTO pizza_sizes (pizza_id, sizes) VALUES (?, ?)"

	id, err = savePizza(ctx, tx, insertPizzaQuery, categoryId, name, description, price, rating, imagePath)
	if err != nil {
		return 0, fmt.Errorf("failed to save pizza: %w", err)
	}

	if err = saveOptions(ctx, tx, insertTypesDoughQuery, id, doughPizzaTypes); err != nil {
		return 0, fmt.Errorf("failed to save type dough: %w", err)
	}

	if err = saveOptions(ctx, tx, insertSizesQuery, id, sizes); err != nil {
		return 0, fmt.Errorf("failed to save sizes: %w", err)
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
	var doughTypes, sizes []int32
	var getPizzaQuery, getTypeDoughQuery, getSizesQuery string

	switch ident.(type) {
	case string:
		getPizzaQuery = "SELECT * FROM pizza WHERE name=?;"
		getTypeDoughQuery = "SELECT dough_type_id FROM pizza_dough_types WHERE pizza_id=(SELECT id FROM pizza where name=?)"
		getSizesQuery = "SELECT sizes FROM pizza_sizes WHERE pizza_id=(SELECT id FROM pizza where name=?)"
	case uint64:
		getPizzaQuery = "SELECT * FROM pizza WHERE id=?;"
		getTypeDoughQuery = "SELECT dough_type_id FROM pizza_dough_types WHERE pizza_id=?"
		getSizesQuery = "SELECT sizes FROM pizza_sizes WHERE pizza_id=?"
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

	if err := getPizzaOptions(ctx, tx, getTypeDoughQuery, ident, &doughTypes); err != nil {
		return nil, err
	}

	if err := getPizzaOptions(ctx, tx, getSizesQuery, ident, &sizes); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	for _, v := range doughTypes {
		pizzaObj.Types = append(pizzaObj.Types, v)
	}

	for _, v := range sizes {
		pizzaObj.Sizes = append(pizzaObj.Sizes, v)
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
func (s *Statement) List(ctx context.Context, indentCategory any, offset, limit uint32) ([]models.Pizza, error) {
	var query string
	args := make([]any, 0)
	pizzas := make([]models.Pizza, 0, limit)

	if reflect.ValueOf(indentCategory).IsZero() {
		query = "SELECT * FROM pizza ORDER BY id LIMIT ? OFFSET ?;"
		args = []any{limit, offset}
	} else {
		switch indentCategory.(type) {
		case string:
			query = "SELECT * FROM pizza WHERE category_id=(SELECT id FROM categories WHERE name=?) ORDER BY id LIMIT ? OFFSET ?;"
		case uint32:
			query = "SELECT * FROM pizza WHERE category_id=? ORDER BY id LIMIT ? OFFSET ?;"
		default:
			return nil, storage.ErrInvalidIdentifier
		}
		args = []any{indentCategory, limit, offset}
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
	getSizesQuery := "SELECT sizes FROM pizza_sizes WHERE pizza_id=?"
	doughType := make([]int32, 0)
	sizes := make([]int32, 0)
	for i := 0; i < len(pizzas); i++ {
		if err := getPizzaOptions(ctx, tx, getTypeDoughQuery, pizzas[i].Id, &doughType); err != nil {
			return nil, err
		}
		for _, doughTypeId := range doughType {
			pizzas[i].Types = append(pizzas[i].Types, doughTypeId)
		}
		if err := getPizzaOptions(ctx, tx, getSizesQuery, pizzas[i].Id, &sizes); err != nil {
			return nil, err
		}
		for _, size := range sizes {
			pizzas[i].Sizes = append(pizzas[i].Sizes, size)
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
	category uint32,
	title string,
	description string,
	types []int32,
	price float32,
	sizes []int32,
	rating uint32,
	imageUrl string,

) error {
	var updateQuery, deleteTypesQuery, insertTypesQuery, deleteSizesQuery, insertSizesQuery string

	updates, args := prepareDataUpdates(map[string]interface{}{
		"category_id": category,
		"name":        title,
		"description": description,
		"price":       price,
		"rating":      rating,
		"image_path":  imageUrl,
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
		insertTypesQuery = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES ((SELECT id FROM pizza WHERE name=?), ?)"
		insertSizesQuery = "INSERT INTO pizza_sizes (pizza_id, sizes) VALUES ((SELECT id FROM pizza WHERE name=?), ?)"
		deleteTypesQuery = "DELETE FROM pizza_dough_types WHERE pizza_id=(SELECT id FROM pizza WHERE name=?)"
		deleteSizesQuery = "DELETE FROM pizza_sizes WHERE pizza_id=(SELECT id FROM pizza WHERE name=?)"
	case uint64:
		updateQuery = "UPDATE pizza SET " + strings.Join(updates, ", ") + " WHERE id=?"
		insertTypesQuery = "INSERT INTO pizza_dough_types (pizza_id, dough_type_id) VALUES (?, ?)"
		insertSizesQuery = "INSERT INTO pizza_sizes (pizza_id, sizes) VALUES (?, ?)"
		deleteTypesQuery = "DELETE FROM pizza_dough_types WHERE pizza_id=?"
		deleteSizesQuery = "DELETE FROM pizza_sizes WHERE pizza_id=?"
	}
	args = append(args, ident)

	if err := updatePizzaTable(ctx, tx, updateQuery, args); err != nil {
		return err
	}

	// Update type dough table
	if err := updatePizzaRelationsTable(ctx, tx, deleteTypesQuery, insertTypesQuery, ident, types); err != nil {
		return err
	}

	// Update sizes table
	if err := updatePizzaRelationsTable(ctx, tx, deleteSizesQuery, insertSizesQuery, ident, sizes); err != nil {
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

	res, err := stmt.ExecContext(ctx, ident)
	if err != nil {
		return false, err
	}

	idRes, err := res.RowsAffected()
	if idRes == 0 {
		return false, sql.ErrNoRows
	}
	if err != nil {
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

	res, err := stmt.ExecContext(ctx, ident)
	if err != nil {
		return false, err
	}

	idRes, err := res.RowsAffected()
	if idRes == 0 {
		return false, sql.ErrNoRows
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// SaveTypeDough save new type dough to the system
func (s *Statement) SaveTypeDough(ctx context.Context, name string) (uint32, error) {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO doughs (name) VALUES (?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, name)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

// GetTypeDough get response with id and name of type dough
func (s *Statement) GetTypeDough(ctx context.Context, id uint32) (name string, err error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT name FROM doughs WHERE id=?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	res := stmt.QueryRowContext(ctx, id)

	if err = res.Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}

// RemoveTypeDough removes type dough by id from the system
func (s *Statement) RemoveTypeDough(ctx context.Context, id uint32) error {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM doughs WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	idRes, err := res.RowsAffected()
	if idRes == 0 {
		return sql.ErrNoRows
	}
	if err != nil {
		return err
	}

	return nil
}
