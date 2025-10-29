package pizzaland

import (
	"context"
	"log/slog"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
)

type Saver interface {
	Save(ctx context.Context, pizzaland *pizzalndv1.PizzaProperties) (pizzaId uint64, err error)
	SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (categoryId uint32, err error)
}

type Getter interface {
	GetById(ctx context.Context, id uint64) (pizza *pizzalndv1.PizzaProperties, err error)
	GetByName(ctx context.Context, name string) (pizza *pizzalndv1.PizzaProperties, err error)
	GetCategoryById(ctx context.Context, id uint64) (category *pizzalndv1.CategoryProperties, err error)
	GetCategoryByName(ctx context.Context, name string) (category *pizzalndv1.CategoryProperties, err error)
	List(ctx context.Context, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	ListCategory(ctx context.Context, name string, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
}

type Remover interface {
	RemoveById(ctx context.Context, id uint64) (success bool, err error)
	RemoveByName(ctx context.Context, name string) (success bool, err error)
	RemoveCategoryById(ctx context.Context, id uint64) (success bool, err error)
	RemoveCategoryByName(ctx context.Context, name string) (success bool, err error)
}

type Updater interface {
	Update(
		ctx context.Context,
		categoryId uint32,
		name string,
		description string,
		typeDough pizzalndv1.TypeDough,
		price float64,
		diameter uint32,
	) (success bool, err error)
	UpdateCategory(ctx context.Context, name string, descriptions string) (success bool, err error)
}

type DomainPizzaLand struct {
	log     *slog.Logger
	saver   Saver
	getter  Getter
	remover Remover
	updater Updater
}

func NewPizzaLand(
	log *slog.Logger,
	saver Saver,
	getter Getter,
	remover Remover,
	updater Updater,
) *DomainPizzaLand {
	return &DomainPizzaLand{
		log:     log,
		saver:   saver,
		getter:  getter,
		remover: remover,
		updater: updater,
	}
}

func (p *DomainPizzaLand) Save(ctx context.Context, pizza *pizzalndv1.PizzaProperties) (pizzaId uint64, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) GetById(ctx context.Context, id uint64) (pizza *pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) GetByName(ctx context.Context, name string) (pizza *pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) List(ctx context.Context, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) CategoryList(ctx context.Context, category string, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) Update(
	ctx context.Context,
	categoryId uint32,
	name, description string,
	typeDough *pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (success bool, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) RemoveById(ctx context.Context, id uint64) (success bool, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) RemoveByName(ctx context.Context, name string) (success bool, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (uint32 uint32, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) GetCategoryById(ctx context.Context, id uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) GetCategoryByName(ctx context.Context, name string) (pizza []*pizzalndv1.PizzaProperties, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) UpdateCategory(ctx context.Context, name, descriptions string) (success bool, err error) {
	panic("implement me")
}

func (p *DomainPizzaLand) RemoveCategoryById(ctx context.Context, id uint32) (success bool, err error) {
	panic("implement me")
}
