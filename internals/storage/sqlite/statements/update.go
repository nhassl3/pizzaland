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

// updatePizzaRelationsTable updates by options (type dough or sizes)
func updatePizzaRelationsTable[T any](ctx context.Context, tx *sql.Tx, deleteQuery, insertQuery string, ident any, opts []T) error {
	if err := deleteOnUpdateOptions(ctx, tx, deleteQuery, ident); err != nil {
		return err
	}

	if len(opts) == 0 {
		return nil
	}

	return insertUpdateOptions(ctx, tx, opts, insertQuery, ident)
}

// insertUpdateOptions inserting new options (type dough or sizes) in storage
func insertUpdateOptions[T any](ctx context.Context, tx *sql.Tx, opts []T, query string, ident any) error {
	if len(opts) == 0 {
		return nil
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare dough types insert: %w", err)
	}
	defer stmt.Close()

	for _, opt := range opts {
		if _, err := stmt.ExecContext(ctx, ident, opt); err != nil {
			return fmt.Errorf("insert dough type %d: %w", opt, err)
		}
	}

	return nil
}

// deleteOnUpdateOptions removes old records about pizza
func deleteOnUpdateOptions(ctx context.Context, tx *sql.Tx, query string, ident any) error {
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, ident); err != nil {
		return fmt.Errorf("failed to execute options delete: %w", err)
	}

	return nil
}
