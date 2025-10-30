package storage

import (
	"errors"
)

var (
	ErrPizzaExists             = errors.New("pizza already exists")
	ErrPizzaNotFound           = errors.New("pizza not found")
	ErrCategoryExists          = errors.New("category already exists")
	ErrCategoryNotFound        = errors.New("category not found")
	ErrNothingToChangePizza    = errors.New("no changes to update pizza record on pizza table")
	ErrNothingToChangeCategory = errors.New("no changes to update category record on categories table")
)
