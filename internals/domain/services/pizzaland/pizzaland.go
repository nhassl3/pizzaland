package pizzaland

import (
	"errors"
	"log/slog"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/lib/logger/sl"
	"github.com/nhassl3/pizzaland/internals/storage"
	"golang.org/x/net/context"
)

const (
	opSave            = "domain.pizzaland.Save"
	opGet             = "domain.pizzaland.Get"
	opList            = "domain.pizzaland.List"
	opUpdate          = "domain.pizzaland.Update"
	opRemove          = "domain.pizzaland.Remove"
	opCategorySave    = "domain.pizzaland.SaveCategory"
	opGetCategory     = "domain.pizzaland.GetCategory"
	opGetCategoryList = "domain.pizzaland.GetCategoryList"
	opUpdateCategory  = "domain.pizzaland.UpdateCategory"
	opRemoveCategory  = "domain.pizzaland.RemoveCategory"
	opSaveTypeDough   = "domain.pizzaland.SaveTypeDough"
	opGetTypeDough    = "domain.pizzaland.GetTypeDough"
	opRemoveTypeDough = "domain.pizzaland.RemoveTypeDough"
)

var ErrPizzaAlreadyExists = errors.New("pizza already exists in the system")
var ErrCategoryAlreadyExists = errors.New("category already exists in the system")
var ErrPizzaNotFound = errors.New("pizza not found in the system")
var ErrCategoryNotFound = errors.New("category not found in the system")
var ErrInvalidIdentifier = errors.New("identifier must be string or uint64")
var ErrListPizzaOutOfRange = errors.New("pizzas list out of range offset")
var ErrTypeDoughAlreadyExists = errors.New("type dough already exists")
var ErrTypeDoughNotFound = errors.New("type dough not found")

type Saver interface {
	Save(ctx context.Context, pizzaland *pizzalndv1.PizzaProperties) (pizzaId uint64, err error)
	SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (categoryId uint32, err error)
	SaveTypeDough(ctx context.Context, name string) (typeDoughId uint32, err error)
}

type Getter interface {
	Get(ctx context.Context, ident any) (pizza *pizzalndv1.PizzaProperties, err error)
	GetCategory(ctx context.Context, ident any) (category *pizzalndv1.CategoryProperties, err error)
	GetTypeDough(ctx context.Context, id uint32) (name string, err error)
	List(ctx context.Context, ident any, offset uint32, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
}

type Remover interface {
	Remove(ctx context.Context, ident any) (success bool, err error)
	RemoveCategory(ctx context.Context, ident any) (success bool, err error)
	RemoveTypeDough(ctx context.Context, id uint32) (success bool, err error)
}

type Updater interface {
	Update(
		ctx context.Context,
		ident any,
		category uint32,
		title string,
		description string,
		types []int32,
		price float32,
		sizes []int32,
		rating uint32,
		imageUrl string,
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

func (p *DomainPizzaLand) List(ctx context.Context, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	return p.list(ctx, 0, offset, limit, opList)
}

func (p *DomainPizzaLand) Update(
	ctx context.Context,
	ident any,
	category uint32,
	title, description string,
	types []int32,
	price float32,
	sizes []int32,
	rating uint32,
	imageUrl string,
) (success bool, err error) {
	log := p.log.With(slog.String("op", opUpdate))

	success, err = p.updater.Update(ctx, ident, category, title, description, types, price, sizes, rating, imageUrl)
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

func (p *DomainPizzaLand) GetCategoryList(ctx context.Context, ident any, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error) {
	return p.list(ctx, ident, offset, limit, opGetCategoryList)
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

func (p *DomainPizzaLand) GetTypeDough(ctx context.Context, id uint32) (name string, err error) {
	log := p.log.With(slog.String("op", opGetTypeDough))

	name, err = p.getter.GetTypeDough(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrTypeDoughNotFound) {
			return "", sl.ErrUpLevel(opGetTypeDough, ErrTypeDoughNotFound)
		}
		log.Error(opGetTypeDough, sl.Err(err))

		return "", sl.ErrUpLevel(opGetTypeDough, err)
	}

	return
}

func (p *DomainPizzaLand) SaveTypeDough(ctx context.Context, name string) (id uint32, err error) {
	log := p.log.With(slog.String("op", opSaveTypeDough))

	id, err = p.saver.SaveTypeDough(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrTypeDoughAlreadyExists) {
			return 0, sl.ErrUpLevel(opSaveTypeDough, ErrTypeDoughAlreadyExists)
		}
		log.Error(opSaveTypeDough, sl.Err(err))

		return 0, sl.ErrUpLevel(opSaveTypeDough, err)
	}

	return
}

func (p *DomainPizzaLand) RemoveTypeDough(ctx context.Context, id uint32) (success bool, err error) {
	log := p.log.With(slog.String("op", opRemoveTypeDough))

	success, err = p.remover.RemoveTypeDough(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrTypeDoughNotFound) {
			return false, sl.ErrUpLevel(opRemoveTypeDough, ErrTypeDoughNotFound)
		}
		log.Error(opRemoveTypeDough, sl.Err(err))

		return false, sl.ErrUpLevel(opRemoveTypeDough, err)
	}

	return
}

func (p *DomainPizzaLand) list(ctx context.Context, ident any, offset, limit uint32, op string) (pizza []*pizzalndv1.PizzaProperties, err error) {
	log := p.log.With(slog.String("op", op))

	pizza, err = p.getter.List(ctx, ident, offset, limit)
	if err != nil {
		if errors.Is(err, storage.ErrPizzaNotFound) {
			return nil, sl.ErrUpLevel(op, ErrPizzaNotFound)
		} else if errors.Is(err, storage.ErrListPizzaOutOfRange) {
			return nil, sl.ErrUpLevel(op, ErrListPizzaOutOfRange)
		}
		log.Error(op, sl.Err(err))

		return nil, sl.ErrUpLevel(op, err)
	}

	return pizza, nil
}
