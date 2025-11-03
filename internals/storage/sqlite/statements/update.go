package statements

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

// prepareDataUpdates prepare data for update event handler
func prepareDataUpdates(params map[string]interface{}) ([]string, []interface{}) {
	updates := make([]string, 0, len(params))
	args := make([]interface{}, 0, len(params))

	addUpdate := func(condition bool, column string, value interface{}) {
		if condition {
			updates = append(updates, column+"=?")
			args = append(args, value)
		}
	}

	for nameField, valueField := range params {
		addUpdate(!reflect.ValueOf(valueField).IsZero(), nameField, valueField)
	}

	return updates, args
}

// updatePizzaTable обновляет основную таблицу пицц
func updatePizzaTable(ctx context.Context, tx *sql.Tx, query string, args []any) error {
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

// updateDoughTypes updates types of the dough for pizza
func updateDoughTypes(ctx context.Context, tx *sql.Tx, deleteQuery, insertQuery string, ident any, typeDough []int32) error {
	if err := deleteDoughTypes(ctx, tx, deleteQuery, ident); err != nil {
		return err
	}

	if len(typeDough) == 0 {
		return nil
	}

	return insertDoughTypes(ctx, tx, typeDough, insertQuery, ident)
}

// deleteDoughTypes remove all dough types for selected identifier of pizza
func deleteDoughTypes(ctx context.Context, tx *sql.Tx, query string, ident any) error {
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

// insertDoughTypes inserting new dough types for pizza
func insertDoughTypes(ctx context.Context, tx *sql.Tx, typeDough []int32, query string, ident any) error {
	if len(typeDough) == 0 {
		return nil
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
