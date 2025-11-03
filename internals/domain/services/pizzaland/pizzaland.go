package pizzaland

import (
	"context"
	"errors"
	"log/slog"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/lib/logger/sl"
	"github.com/nhassl3/pizzaland/internals/storage"
)

const (
	opSave           = "domain.pizzaland.Save"
	opGet            = "domain.pizzaland.Get"
	opList           = "domain.pizzaland.List"
	opUpdate         = "domain.pizzaland.Update"
	opRemove         = "domain.pizzaland.Remove"
	opCategorySave   = "domain.pizzaland.SaveCategory"
	opGetCategory    = "domain.pizzaland.GetCategory"
	opUpdateCategory = "domain.pizzaland.UpdateCategory"
	opRemoveCategory = "domain.pizzaland.RemoveCategory"
)

var (
	ErrPizzaAlreadyExists    = errors.New("pizza already exists in the system")
	ErrCategoryAlreadyExists = errors.New("category already exists in the system")
	ErrPizzaNotFound         = errors.New("pizza not found in the system")
	ErrCategoryNotFound      = errors.New("category not found in the system")
	ErrInvalidIdentifier     = errors.New("identifier must be string or uint64")
	ErrListPizzaOutOfRange   = errors.New("pizzas list out of range offset")
)

type Saver interface {
	Save(ctx context.Context, pizzaland *pizzalndv1.PizzaProperties) (pizzaId uint64, err error)
	SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (categoryId uint32, err error)
}

type Getter interface {
	Get(ctx context.Context, ident any) (pizza *pizzalndv1.PizzaProperties, err error)
	GetCategory(ctx context.Context, ident any) (category *pizzalndv1.CategoryProperties, err error)
	List(ctx context.Context, ident any, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
}

type Remover interface {
	Remove(ctx context.Context, ident any) (success bool, err error)
	RemoveCategory(ctx context.Context, ident any) (success bool, err error)
}

type Updater interface {
	Update(
		ctx context.Context,
		ident any,
		categoryId uint32,
		name string,
		description string,
		typeDough []pizzalndv1.TypeDough,
		price float32,
		diameter uint32,
	) (success bool, err error)
	UpdateCategoryById(ctx context.Context, id uint32, name string, descriptions string) (success bool, err error)
	UpdateCategoryByName(ctx context.Context, name string, descriptions string) (success bool, err error)
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
	log := p.log.With(slog.String("op", opSave))

	pizzaId, err = p.saver.Save(ctx, pizza)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaExists) {
			return 0, sl.ErrUpLevel(opSave, ErrPizzaAlreadyExists)
		}
		log.Error(opSave, sl.Err(err))

		return 0, sl.ErrUpLevel(opSave, err)
	}

	return pizzaId, nil
}

func (p *DomainPizzaLand) Get(ctx context.Context, ident any) (pizza *pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", opGet))

	pizza, err = p.getter.Get(ctx, ident)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(opGet, ErrPizzaNotFound)
		} else if errors.Is(err, storage.ErrInvalidIdentifier) {
			return nil, sl.ErrUpLevel(opGet, ErrInvalidIdentifier)
		}
		log.Error(opGet, sl.Err(err))

		return nil, sl.ErrUpLevel(opGet, err)
	}

	return pizza, nil
}

func (p *DomainPizzaLand) List(ctx context.Context, ident any, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", opList))

	pizza, err = p.getter.List(ctx, ident, offset, limit)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(opList, ErrPizzaNotFound)
		} else if errors.Is(err, storage.ErrListPizzaOutOfRange) {
			return nil, sl.ErrUpLevel(opList, ErrListPizzaOutOfRange)
		}
		log.Error(opList, sl.Err(err))

		return nil, sl.ErrUpLevel(opList, err)
	}

	return pizza, nil
}

func (p *DomainPizzaLand) Update(
	ctx context.Context,
	ident any,
	categoryId uint32,
	name, description string,
	typeDough []pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (success bool, err error) {
	log := p.log.With(slog.String("op", opUpdate))

	success, err = p.updater.Update(ctx, ident, categoryId, name, description, typeDough, price, diameter)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return false, sl.ErrUpLevel(opUpdate, ErrPizzaNotFound)
		} else if errors.Is(err, storage.ErrPizzaExists) {
			return false, sl.ErrUpLevel(opUpdate, ErrPizzaAlreadyExists)
		} else if errors.Is(err, storage.ErrInvalidIdentifier) {
			return false, sl.ErrUpLevel(opUpdate, ErrInvalidIdentifier)
		}
		log.Error(opUpdate, sl.Err(err))

		return false, sl.ErrUpLevel(opUpdate, err)
	}

	return
}

func (p *DomainPizzaLand) Remove(ctx context.Context, ident any) (success bool, err error) {
	log := p.log.With(slog.String("op", opRemove))

	success, err = p.remover.Remove(ctx, ident)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return false, sl.ErrUpLevel(opRemove, ErrPizzaNotFound)
		}
		log.Error(opRemove, sl.Err(err))

		return false, sl.ErrUpLevel(opRemove, err)
	}

	return
}

func (p *DomainPizzaLand) SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (uint32 uint32, err error) {
	log := p.log.With(slog.String("op", opCategorySave))

	categoryId, err := p.saver.SaveCategory(ctx, category)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryExists) {
			return 0, sl.ErrUpLevel(opCategorySave, ErrCategoryAlreadyExists)
		}
		log.Error(opCategorySave, sl.Err(err))
		return 0, sl.ErrUpLevel(opCategorySave, err)
	}

	return categoryId, nil
}

func (p *DomainPizzaLand) GetCategory(ctx context.Context, ident any) (category *pizzalndv1.CategoryProperties, err error) {
	log := p.log.With(slog.String("op", opGetCategory))

	category, err = p.getter.GetCategory(ctx, ident)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			return nil, sl.ErrUpLevel(opGetCategory, ErrCategoryNotFound)
		}
		log.Error(opGetCategory, sl.Err(err))

		return nil, sl.ErrUpLevel(opGetCategory, err)
	}

	return
}

func (p *DomainPizzaLand) UpdateCategory(ctx context.Context, id uint32, name, description string) (success bool, err error) {
	log := p.log.With(slog.String("op", opUpdateCategory))

	success, err = p.updater.UpdateCategoryById(ctx, id, name, description)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			return false, sl.ErrUpLevel(opUpdateCategory, ErrCategoryNotFound)
		}
		log.Error(opUpdateCategory, sl.Err(err))

		return false, sl.ErrUpLevel(opUpdateCategory, err)
	}

	return
}

func (p *DomainPizzaLand) RemoveCategory(ctx context.Context, ident any) (success bool, err error) {
	log := p.log.With(slog.String("op", opRemoveCategory))

	success, err = p.remover.RemoveCategory(ctx, ident)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			return false, sl.ErrUpLevel(opRemoveCategory, ErrCategoryNotFound)
		}
		log.Error(opRemoveCategory, sl.Err(err))

		return false, sl.ErrUpLevel(opRemoveCategory, err)
	}

	return
}
