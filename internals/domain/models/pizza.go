package models

import (
	"time"
)

// Pizza Naming fields of the structure must be equal with naming in api of the project
type Pizza struct {
	Id          uint64
	Category    uint32
	Title       string
	Description string
	Types       []int32
	Price       float64
	Sizes       []int32
	Rating      uint32
	ImageUrl    string
	CreatedAt   time.Time
}
