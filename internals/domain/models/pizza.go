package models

type Pizza struct {
	PizzaId     uint64
	CategoryId  uint32
	Name        string
	Description string
	TypeDough   int64
	Price       float64
	Diameter    uint32
}
