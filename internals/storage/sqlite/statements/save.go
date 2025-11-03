package statements

import (
	"context"
	"database/sql"
	"fmt"
)

func savePizza(
	ctx context.Context,
	tx *sql.Tx,
	query string,
	args ...any,
) (uint64, error) {
	insertPizza, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer insertPizza.Close()

	res, err := insertPizza.ExecContext(ctx, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to insert pizza record: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return uint64(id), nil
}

func saveTypeDough(
	ctx context.Context,
	tx *sql.Tx,
	query string,
	pizzaId uint64,
	doughPizzaTypes []int32,
) error {
	if len(doughPizzaTypes) == 0 {
		return nil
	}

	stmt, err := tx.PrepareContext(ctx, query)
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
