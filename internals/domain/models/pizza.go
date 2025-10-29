package models

type Pizza struct {
	ID          uint64
	CategoryId  uint32
	Name        string
	Description string
	Price       float64
	TypeDough   TypeDough
	Diameter    uint32
}

type TypeDough struct {
	ID   uint32
	Name string
}
