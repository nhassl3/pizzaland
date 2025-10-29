package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, pizza *pizzalndv1.PizzaProperties) (pizzaId uint64, err error) {
	panic("implement me")
}

func (s *Storage) SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (categoryId uint32, err error) {
	panic("implement me")
}

func (s *Storage) GetById(ctx context.Context, id uint64) (pizza *pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (s *Storage) GetByName(ctx context.Context, name string) (pizza *pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (s *Storage) GetCategoryById(ctx context.Context, id uint64) (category *pizzalndv1.CategoryProperties, err error) {
	panic("implement me")
}

func (s *Storage) GetCategoryByName(ctx context.Context, name string) (category *pizzalndv1.CategoryProperties, err error) {
	panic("implement me")
}

func (s *Storage) List(ctx context.Context, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (s *Storage) ListCategory(ctx context.Context, name string, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (s *Storage) RemoveById(ctx context.Context, id uint64) (success bool, err error) {
	panic("implement me")
}

func (s *Storage) RemoveByName(ctx context.Context, name string) (success bool, err error) {
	panic("implement me")
}

func (s *Storage) RemoveCategoryById(ctx context.Context, id uint64) (success bool, err error) {
	panic("implement me")
}

func (s *Storage) RemoveCategoryByName(ctx context.Context, name string) (success bool, err error) {
	panic("implement me")
}

func (s *Storage) Update(
	ctx context.Context,
	categoryId uint32,
	name string,
	description string,
	typeDough pizzalndv1.TypeDough,
	price float64,
	diameter uint32,
) (success bool, err error) {
	panic("implement me")
}

func (s *Storage) UpdateCategory(ctx context.Context, name string, descriptions string) (success bool, err error) {
	panic("implement me")
}
