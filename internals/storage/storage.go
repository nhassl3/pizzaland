package storage

import "errors"

var (
	ErrPizzaExists      = errors.New("pizza already exists")
	ErrPizzaNotFound    = errors.New("pizza not found")
	ErrCategoryExists   = errors.New("category already exists")
	ErrCategoryNotFound = errors.New("category not found")
)
