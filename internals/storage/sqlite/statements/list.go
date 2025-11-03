package statements

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nhassl3/pizzaland/internals/domain/models"
)

func getPizzasSlice(ctx context.Context, tx *sql.Tx, query string, pizzaObj *[]models.Pizza, args []any) error {
	stmt, err := tx.PrepareContext(
		ctx,
		query,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}

	var model models.Pizza
	for rows.Next() {
		if err = rows.Scan(
			&model.PizzaId,
			&model.CategoryId,
			&model.Name,
			&model.Description,
			&model.Price,
			&model.Diameter,
			&model.CreatedAt,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			return err
		}
		*pizzaObj = append(*pizzaObj, model)
	}

	return nil
}
