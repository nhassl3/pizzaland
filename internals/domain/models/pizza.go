package models

import (
	"time"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
)

type Pizza struct {
	PizzaId     uint64
	CategoryId  uint32
	Name        string
	Description string
	TypeDough   []pizzalndv1.TypeDough
	Price       float64
	Diameter    uint32
	CreatedAt   time.Time
}
