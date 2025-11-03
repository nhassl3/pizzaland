package statements

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/nhassl3/pizzaland/internals/domain/models"
	"github.com/nhassl3/pizzaland/internals/storage"
)

func getPizza(ctx context.Context, tx *sql.Tx, query string, ident any, pizzaObj *models.Pizza) error {
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare request: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, ident)

	if err = row.Scan(
		&pizzaObj.PizzaId,
		&pizzaObj.CategoryId,
		&pizzaObj.Name,
		&pizzaObj.Description,
		&pizzaObj.Price,
		&pizzaObj.Diameter,
		&pizzaObj.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrPizzaNotFound
		}
		return err
	}

	return nil
}

func getTypeDough(
	ctx context.Context,
	tx *sql.Tx,
	query string,
	ident any,
	doughTypes *[]int32,
) error {
	var doughType int32

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, ident)
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.Scan(&doughType); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return err
		}
		*doughTypes = append(*doughTypes, doughType)
	}

	return nil
}
