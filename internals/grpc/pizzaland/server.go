package pizzaland

import (
	"context"
	"errors"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/domain/services/pizzaland"
	"github.com/nhassl3/pizzaland/internals/lib/reflection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	NoIdentifier           = "None of several arguments were provided"
	UnknownNameOrId        = "An unknown Title or ID of the pizza was given"
	PizzaAlreadyExists     = "Pizza already exists"
	CategoryAlreadyExists  = "Category already exists"
	PizzaNotFound          = "PizzaNotFound"
	CategoryNotFound       = "Category not found"
	ChangeCapacity         = "Change capacity of offset"
	TypeDoughAlreadyExists = "Type dough already exists"
	TypeDoughNotFound      = "Type dough not found"
)

type PizzaLand interface {
	Save(ctx context.Context, pizza *pizzalndv1.PizzaProperties) (pizzaId uint64, err error)
	Get(ctx context.Context, ident any) (pizza *pizzalndv1.PizzaProperties, err error)
	List(ctx context.Context, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	Update(
		ctx context.Context,
		ident any,
		category uint32,
		title, description string,
		types []int32,
		price float32,
		sizes []int32,
		rating uint32,
		imageUrl string,
	) (success bool, err error)
	Remove(ctx context.Context, ident any) (success bool, err error)
	SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (uint32 uint32, err error)
	GetCategory(ctx context.Context, ident any) (category *pizzalndv1.CategoryProperties, err error)
	GetCategoryList(ctx context.Context, ident any, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	UpdateCategory(ctx context.Context, id uint32, name, descriptions string) (success bool, err error)
	RemoveCategory(ctx context.Context, ident any) (success bool, err error)
	GetTypeDough(ctx context.Context, id uint32) (name string, err error)
	SaveTypeDough(ctx context.Context, name string) (id uint32, err error)
	RemoveTypeDough(ctx context.Context, id uint32) (success bool, err error)
}

type ServerAPI struct {
	pizzalndv1.UnimplementedPizzaLandServer
	pizzaLand PizzaLand
}

func Register(gRPCServer *grpc.Server, pizzaLand PizzaLand) {
	pizzalndv1.RegisterPizzaLandServer(gRPCServer, &ServerAPI{pizzaLand: pizzaLand})
}

func (api *ServerAPI) Save(ctx context.Context, in *pizzalndv1.SaveRequest) (*pizzalndv1.SaveResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pizzaId, err := api.pizzaLand.Save(ctx, in.GetPizza())
	if err != nil {
		if errors.Is(err, pizzaland.ErrPizzaAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, PizzaAlreadyExists)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.SaveResponse{PizzaId: pizzaId}, nil
}

func (api *ServerAPI) Get(ctx context.Context, in *pizzalndv1.GetRequest) (*pizzalndv1.GetResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		pizza *pizzalndv1.PizzaProperties
		err   error
	)

	switch v := in.GetIdentifier().(type) {
	case *pizzalndv1.GetRequest_PizzaId:
		pizza, err = api.pizzaLand.Get(ctx, v.PizzaId)
	case *pizzalndv1.GetRequest_PizzaName:
		pizza, err = api.pizzaLand.Get(ctx, v.PizzaName)
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		if errors.Is(err, pizzaland.ErrInvalidIdentifier) {
			return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
		} else if errors.Is(err, pizzaland.ErrPizzaNotFound) {
			return nil, status.Error(codes.NotFound, PizzaNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.GetResponse{Pizza: pizza}, nil
}

func (api *ServerAPI) List(ctx context.Context, in *pizzalndv1.ListRequest) (*pizzalndv1.ListResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pizza, err := api.pizzaLand.List(ctx, in.GetOffset(), in.GetLimit())

	if err != nil {
		if errors.Is(err, pizzaland.ErrListPizzaOutOfRange) {
			return nil, status.Error(codes.InvalidArgument, ChangeCapacity)
		} else if errors.Is(err, pizzaland.ErrPizzaNotFound) {
			return nil, status.Error(codes.NotFound, PizzaNotFound)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.ListResponse{Pizza: pizza}, nil
}

func (api *ServerAPI) Update(ctx context.Context, in *pizzalndv1.UpdateRequest) (*pizzalndv1.UpdateResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if reflection.AllFieldsIsNil(in) {
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	}

	var ident any
	switch v := in.GetIdentifier().(type) {
	case *pizzalndv1.UpdateRequest_Id:
		ident = v.Id
	case *pizzalndv1.UpdateRequest_SourceName:
		ident = v.SourceName
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	var (
		category    = in.GetCategory().GetValue()
		title       = in.GetTitle().GetValue()
		description = in.GetDescription().GetValue()
		price       = in.GetPrice().GetValue()
		sizes       = in.GetSizes()
		types       = in.GetTypes()
		rating      = in.GetRating().GetValue()
		imageUrl    = in.GetImageUrl().GetValue()
	)

	success, err := api.pizzaLand.Update(ctx, ident, category, title, description, types, price, sizes, rating, imageUrl)
	if err != nil {
		if errors.Is(err, pizzaland.ErrPizzaNotFound) {
			return nil, status.Error(codes.NotFound, PizzaNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.UpdateResponse{Success: success}, nil
}

func (api *ServerAPI) Remove(ctx context.Context, in *pizzalndv1.RemoveRequest) (*pizzalndv1.RemoveResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		success bool
		err     error
	)

	switch v := in.GetIdentifier().(type) {
	case *pizzalndv1.RemoveRequest_PizzaId:
		success, err = api.pizzaLand.Remove(ctx, v.PizzaId)
	case *pizzalndv1.RemoveRequest_PizzaName:
		success, err = api.pizzaLand.Remove(ctx, v.PizzaName)
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		if errors.Is(err, pizzaland.ErrPizzaNotFound) {
			return nil, status.Error(codes.NotFound, PizzaNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.RemoveResponse{Success: success}, nil
}

func (api *ServerAPI) SaveCategory(ctx context.Context, in *pizzalndv1.SaveCategoryRequest) (*pizzalndv1.SaveCategoryResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	categoryId, err := api.pizzaLand.SaveCategory(ctx, in.GetCategory())
	if err != nil {
		if errors.Is(err, pizzaland.ErrCategoryAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, CategoryAlreadyExists)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.SaveCategoryResponse{CategoryId: categoryId}, nil
}

func (api *ServerAPI) GetCategory(ctx context.Context, in *pizzalndv1.GetCategoryRequest) (*pizzalndv1.GetCategoryResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		category *pizzalndv1.CategoryProperties
		err      error
	)

	switch v := in.GetIdentifier().(type) {
	case *pizzalndv1.GetCategoryRequest_CategoryId:
		category, err = api.pizzaLand.GetCategory(ctx, v.CategoryId)
	case *pizzalndv1.GetCategoryRequest_CategoryName:
		category, err = api.pizzaLand.GetCategory(ctx, v.CategoryName)
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		if errors.Is(err, pizzaland.ErrInvalidIdentifier) {
			return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
		} else if errors.Is(err, pizzaland.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, CategoryNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.GetCategoryResponse{Category: category}, nil
}

func (api *ServerAPI) GetCategoryList(ctx context.Context, in *pizzalndv1.GetCategoryListRequest) (*pizzalndv1.GetCategoryListResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		pizza []*pizzalndv1.PizzaProperties
		err   error
	)

	switch v := in.GetIdentifier().(type) {
	case *pizzalndv1.GetCategoryListRequest_CategoryId:
		pizza, err = api.pizzaLand.GetCategoryList(ctx, v.CategoryId, in.GetOffset(), in.GetLimit())
	case *pizzalndv1.GetCategoryListRequest_CategoryName:
		pizza, err = api.pizzaLand.GetCategoryList(ctx, v.CategoryName, in.GetOffset(), in.GetLimit())
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		if errors.Is(err, pizzaland.ErrInvalidIdentifier) {
			return nil, status.Error(codes.InvalidArgument, NoIdentifier)
		} else if errors.Is(err, pizzaland.ErrListPizzaOutOfRange) {
			return nil, status.Error(codes.InvalidArgument, ChangeCapacity)
		} else if errors.Is(err, pizzaland.ErrPizzaNotFound) {
			return nil, status.Error(codes.NotFound, PizzaNotFound)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.GetCategoryListResponse{
		Pizza: &pizzalndv1.ListResponse{
			Pizza: pizza,
		},
	}, nil
}

func (api *ServerAPI) UpdateCategory(ctx context.Context, in *pizzalndv1.UpdateCategoryRequest) (*pizzalndv1.UpdateCategoryResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if reflection.AllFieldsIsNil(in) {
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	}

	var (
		id          = in.GetCategoryId()
		name        = in.GetName().GetValue()
		description = in.GetDescription().GetValue()
	)

	success, err := api.pizzaLand.UpdateCategory(ctx, id, name, description)
	if err != nil {
		if errors.Is(err, pizzaland.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, CategoryNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.UpdateCategoryResponse{Success: success}, nil
}

func (api *ServerAPI) RemoveCategory(ctx context.Context, in *pizzalndv1.RemoveCategoryRequest) (*pizzalndv1.RemoveCategoryResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		success bool
		err     error
	)

	switch v := in.GetIdentifier().(type) {
	case *pizzalndv1.RemoveCategoryRequest_CategoryId:
		success, err = api.pizzaLand.RemoveCategory(ctx, v.CategoryId)
	case *pizzalndv1.RemoveCategoryRequest_CategoryName:
		success, err = api.pizzaLand.RemoveCategory(ctx, v.CategoryName)
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		if errors.Is(err, pizzaland.ErrCategoryNotFound) {
			return nil, status.Error(codes.NotFound, CategoryNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.RemoveCategoryResponse{Success: success}, nil
}

func (api *ServerAPI) SaveTypeDough(ctx context.Context, in *pizzalndv1.SaveTypeDoughRequest) (*pizzalndv1.SaveTypeDoughResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := api.pizzaLand.SaveTypeDough(ctx, in.GetName())
	if err != nil {
		if errors.Is(err, pizzaland.ErrTypeDoughAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, TypeDoughAlreadyExists)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.SaveTypeDoughResponse{Id: id}, nil
}

func (api *ServerAPI) GetTypeDough(ctx context.Context, in *pizzalndv1.GetTypeDoughRequest) (*pizzalndv1.GetTypeDoughResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	name, err := api.pizzaLand.GetTypeDough(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, pizzaland.ErrTypeDoughNotFound) {
			return nil, status.Error(codes.NotFound, TypeDoughNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.GetTypeDoughResponse{Name: name}, nil
}

func (api *ServerAPI) RemoveTypeDough(ctx context.Context, in *pizzalndv1.RemoveTypeDoughRequest) (*pizzalndv1.RemoveTypeDoughResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	success, err := api.pizzaLand.RemoveTypeDough(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, pizzaland.ErrTypeDoughNotFound) {
			return nil, status.Error(codes.NotFound, TypeDoughNotFound)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.RemoveTypeDoughResponse{Success: success}, nil
}
