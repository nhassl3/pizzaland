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
	opSave                 = "domain.pizzaland.Save"
	opGetById              = "domain.pizzaland.GetById"
	opGetByName            = "domain.pizzaland.GetByName"
	opList                 = "domain.pizzaland.List"
	opCategoryList         = "domain.pizzaland.CategoryList"
	opCategoryListById     = "domain.pizzaland.CategoryListById"
	opUpdate               = "domain.pizzaland.Update"
	opRemoveById           = "domain.pizzaland.RemoveById"
	opRemoveByName         = "domain.pizzaland.RemoveByName"
	opCategorySave         = "domain.pizzaland.SaveCategory"
	opCategoryGetById      = "domain.pizzaland.GetCategoryById"
	opCategoryGetByName    = "domain.pizzaland.GetCategoryByName"
	opUpdateCategory       = "domain.pizzaland.UpdateCategory"
	opRemoveCategoryById   = "domain.pizzaland.RemoveCategoryById"
	opRemoveCategoryByName = "domain.pizzaland.RemoveCategoryByName"
)

var (
	ErrPizzaAlreadyExists    = errors.New("pizza already exists in the system")
	ErrCategoryAlreadyExists = errors.New("category already exists in the system")
	ErrPizzaNotFound         = errors.New("pizza not found in the system")
	ErrCategoryNotFound      = errors.New("category not found in the system")
)

type Saver interface {
	Save(ctx context.Context, pizzaland *pizzalndv1.PizzaProperties) (pizzaId uint64, err error)
	SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (categoryId uint32, err error)
}

type Getter interface {
	GetById(ctx context.Context, id uint64) (pizza *pizzalndv1.PizzaProperties, err error)
	GetByName(ctx context.Context, name string) (pizza *pizzalndv1.PizzaProperties, err error)
	GetCategoryById(ctx context.Context, id uint32) (category *pizzalndv1.CategoryProperties, err error)
	GetCategoryByName(ctx context.Context, name string) (category *pizzalndv1.CategoryProperties, err error)
	List(ctx context.Context, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	ListCategoryByName(ctx context.Context, categoryName string, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	ListCategoryById(ctx context.Context, categoryId, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
}

type Remover interface {
	RemoveById(ctx context.Context, id uint64) (success bool, err error)
	RemoveByName(ctx context.Context, name string) (success bool, err error)
	RemoveCategoryById(ctx context.Context, id uint32) (success bool, err error)
	RemoveCategoryByName(ctx context.Context, name string) (success bool, err error)
}

type Updater interface {
	UpdateById(
		ctx context.Context,
		id uint64,
		categoryId uint32,
		name string,
		description string,
		typeDough *pizzalndv1.TypeDough,
		price float32,
		diameter uint32,
	) (success bool, err error)
	UpdateByName(
		ctx context.Context,
		category uint32,
		name string,
		description string,
		typeDough *pizzalndv1.TypeDough,
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
			return 0, sl.ErrUpLevel(opSave, ErrPizzaAlreadyExists.Error())
		}
		log.Error(opSave, sl.Err(err))

		return 0, sl.ErrUpLevel(opSave, err.Error())
	}

	return pizzaId, nil
}

func (p *DomainPizzaLand) GetById(ctx context.Context, id uint64) (pizza *pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", opGetById))

	pizza, err = p.getter.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(opGetById, ErrPizzaNotFound.Error())
		}
		log.Error(opGetById, sl.Err(err))

		return nil, sl.ErrUpLevel(opGetById, err.Error())
	}

	return pizza, nil
}

func (p *DomainPizzaLand) GetByName(ctx context.Context, name string) (pizza *pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", opGetByName))

	pizza, err = p.getter.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(opGetByName, ErrPizzaNotFound.Error())
		}
		log.Error(opGetByName, sl.Err(err))

		return nil, sl.ErrUpLevel(opGetByName, err.Error())
	}

	return pizza, nil
}

func (p *DomainPizzaLand) List(ctx context.Context, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", opList))

	pizza, err = p.getter.List(ctx, offset, limit)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(opList, ErrPizzaNotFound.Error())
		}
		log.Error(opList, sl.Err(err))

		return nil, sl.ErrUpLevel(opList, err.Error())
	}

	return pizza, nil
}

func (p *DomainPizzaLand) CategoryList(ctx context.Context, category string, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", opCategoryList))

	pizza, err = p.getter.ListCategoryByName(ctx, category, offset, limit)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(opCategoryList, ErrPizzaNotFound.Error())
		}
		log.Error(opCategoryList, sl.Err(err))

		return nil, sl.ErrUpLevel(opCategoryList, err.Error())
	}

	return pizza, nil
}

func (p *DomainPizzaLand) CategoryListById(ctx context.Context, id uint32, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", opCategoryListById))

	pizza, err = p.getter.ListCategoryById(ctx, id, offset, limit)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(opCategoryListById, ErrPizzaNotFound.Error())
		}
		log.Error(opCategoryListById, sl.Err(err))

		return nil, sl.ErrUpLevel(opCategoryListById, err.Error())
	}

	return pizza, nil
}

func (p *DomainPizzaLand) Update(
	ctx context.Context,
	id uint64,
	categoryId uint32,
	name, description string,
	typeDough *pizzalndv1.TypeDough,
	price float32,
	diameter uint32,
) (success bool, err error) {
	log := p.log.With(slog.String("op", opUpdate))

	success, err = p.updater.UpdateById(ctx, id, categoryId, name, description, typeDough, price, diameter)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return false, sl.ErrUpLevel(opUpdate, ErrPizzaNotFound.Error())
		}
		log.Error(opUpdate, sl.Err(err))

		return false, sl.ErrUpLevel(opUpdate, err.Error())
	}

	return
}

func (p *DomainPizzaLand) RemoveById(ctx context.Context, id uint64) (success bool, err error) {
	log := p.log.With(slog.String("op", opRemoveById))

	success, err = p.remover.RemoveById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return false, sl.ErrUpLevel(opRemoveById, ErrPizzaNotFound.Error())
		}
		log.Error(opRemoveById, sl.Err(err))

		return false, sl.ErrUpLevel(opRemoveById, err.Error())
	}

	return
}

func (p *DomainPizzaLand) RemoveByName(ctx context.Context, name string) (success bool, err error) {
	log := p.log.With(slog.String("op", opRemoveByName))

	success, err = p.remover.RemoveByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return false, sl.ErrUpLevel(opRemoveByName, ErrPizzaNotFound.Error())
		}
		log.Error(opRemoveByName, sl.Err(err))

		return false, sl.ErrUpLevel(opRemoveByName, err.Error())
	}

	return
}

func (p *DomainPizzaLand) SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (uint32 uint32, err error) {
	log := p.log.With(slog.String("op", opCategorySave))

	categoryId, err := p.saver.SaveCategory(ctx, category)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryExists) {
			return 0, sl.ErrUpLevel(opCategorySave, ErrCategoryAlreadyExists.Error())
		}
		log.Error(opCategorySave, sl.Err(err))
		return 0, sl.ErrUpLevel(opCategorySave, err.Error())
	}

	return categoryId, nil
}

func (p *DomainPizzaLand) GetCategoryById(ctx context.Context, id uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	return p.CategoryListById(ctx, id, 0, 12)
}

func (p *DomainPizzaLand) GetCategoryByName(ctx context.Context, name string) (pizza []*pizzalndv1.PizzaProperties, err error) {
	return p.CategoryList(ctx, name, 0, 12)
}

func (p *DomainPizzaLand) UpdateCategory(ctx context.Context, id uint32, name, description string) (success bool, err error) {
	log := p.log.With(slog.String("op", opUpdateCategory))

	success, err = p.updater.UpdateCategoryById(ctx, id, name, description)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			return false, sl.ErrUpLevel(opUpdateCategory, ErrCategoryNotFound.Error())
		}
		log.Error(opUpdateCategory, sl.Err(err))

		return false, sl.ErrUpLevel(opUpdateCategory, err.Error())
	}

	return
}

func (p *DomainPizzaLand) RemoveCategoryById(ctx context.Context, id uint32) (success bool, err error) {
	log := p.log.With(slog.String("op", opRemoveCategoryById))

	success, err = p.remover.RemoveCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			return false, sl.ErrUpLevel(opRemoveCategoryById, ErrCategoryNotFound.Error())
		}
		log.Error(opRemoveCategoryById, sl.Err(err))

		return false, sl.ErrUpLevel(opRemoveCategoryById, err.Error())
	}

	return
}

func (p *DomainPizzaLand) RemoveCategoryByName(ctx context.Context, name string) (success bool, err error) {
	log := p.log.With(slog.String("op", opRemoveCategoryByName))

	success, err = p.remover.RemoveCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrCategoryNotFound) {
			return false, sl.ErrUpLevel(opRemoveCategoryByName, ErrCategoryNotFound.Error())
		}
		log.Error(opRemoveCategoryByName, sl.Err(err))

		return false, sl.ErrUpLevel(opRemoveCategoryByName, err.Error())
	}

	return
}
