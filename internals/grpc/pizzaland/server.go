package pizzaland

import (
	"context"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/lib/reflection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	NoIdentifier    = "None of several arguments were provided"
	UnknownNameOrId = "An unknown Name or ID of the pizza was given"
)

type PizzaLand interface {
	Save(ctx context.Context, pizza *pizzalndv1.PizzaProperties) (pizzaId uint64, err error)
	GetById(ctx context.Context, id uint64) (pizza *pizzalndv1.PizzaProperties, err error)
	GetByName(ctx context.Context, name string) (pizza *pizzalndv1.PizzaProperties, err error)
	List(ctx context.Context, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	CategoryList(ctx context.Context, category string, offset, limit uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	Update(
		ctx context.Context,
		id uint64,
		categoryId uint32,
		name, description string,
		typeDough *pizzalndv1.TypeDough,
		price float32,
		diameter uint32,
	) (success bool, err error)
	RemoveById(ctx context.Context, id uint64) (success bool, err error)
	RemoveByName(ctx context.Context, name string) (success bool, err error)
	SaveCategory(ctx context.Context, category *pizzalndv1.CategoryProperties) (uint32 uint32, err error)
	GetCategoryById(ctx context.Context, id uint32) (pizza []*pizzalndv1.PizzaProperties, err error)
	GetCategoryByName(ctx context.Context, name string) (pizza []*pizzalndv1.PizzaProperties, err error)
	UpdateCategory(ctx context.Context, id uint32, name, descriptions string) (success bool, err error)
	RemoveCategoryById(ctx context.Context, id uint32) (success bool, err error)
	RemoveCategoryByName(ctx context.Context, name string) (success bool, err error)
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
		// TODO: no internal error check

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
		pizza, err = api.pizzaLand.GetById(ctx, v.PizzaId)
	case *pizzalndv1.GetRequest_PizzaName:
		pizza, err = api.pizzaLand.GetByName(ctx, v.PizzaName)
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		// TODO: no internal error check
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.GetResponse{Pizza: pizza}, nil
}

func (api *ServerAPI) List(ctx context.Context, in *pizzalndv1.ListRequest) (*pizzalndv1.ListResponse, error) {
	pizza := make([]*pizzalndv1.PizzaProperties, 0, in.GetOffset()) // capacity is offset

	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var err error
	if in.GetCategoryName() != nil {
		pizza, err = api.pizzaLand.CategoryList(ctx, in.GetCategoryName().GetValue(), in.GetOffset(), in.GetLimit())
	} else {
		pizza, err = api.pizzaLand.List(ctx, in.GetOffset(), in.GetLimit())
	}

	if err != nil {
		// TODO: no internal check
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

	var (
		id          = in.GetId()
		categoryId  = in.GetCategoryId().GetValue()
		name        = in.GetName().GetValue()
		description = in.GetDescription().GetValue()
		price       = in.GetPrice().GetValue()
		diameter    = in.GetDiameter().GetValue()
		typeDough   = in.GetTypeDough().Enum()
	)

	success, err := api.pizzaLand.Update(ctx, id, categoryId, name, description, typeDough, price, diameter)
	if err != nil {
		// TODO: no internal error check
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
		success, err = api.pizzaLand.RemoveById(ctx, v.PizzaId)
	case *pizzalndv1.RemoveRequest_PizzaName:
		success, err = api.pizzaLand.RemoveByName(ctx, v.PizzaName)
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		// TODO: no internal error check
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
		// TODO: implement no internal error check
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.SaveCategoryResponse{CategoryId: categoryId}, nil
}

func (api *ServerAPI) GetCategory(ctx context.Context, in *pizzalndv1.GetCategoryRequest) (*pizzalndv1.GetCategoryResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		pizza *pizzalndv1.ListResponse
		err   error
	)

	switch v := in.GetIdentifier().(type) {
	case *pizzalndv1.GetCategoryRequest_CategoryId:
		pizza, err = api.List(ctx, &pizzalndv1.ListRequest{
			CategoryId: wrapperspb.UInt32(v.CategoryId),
			Offset:     0,
			Limit:      12,
		})
	case *pizzalndv1.GetCategoryRequest_CategoryName:
		pizza, err = api.List(ctx, &pizzalndv1.ListRequest{
			CategoryName: wrapperspb.String(v.CategoryName),
			Offset:       0,
			Limit:        12,
		})
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		// TODO: implement no internal error check
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.GetCategoryResponse{Pizza: pizza}, nil
}

func (api *ServerAPI) UpdateCategory(ctx context.Context, in *pizzalndv1.UpdateCategoryRequest) (*pizzalndv1.UpdateCategoryResponse, error) {
	if err := in.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if reflection.AllFieldsIsNil(in.ProtoReflect()) {
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	}

	var (
		id          = in.GetCategoryId()
		name        = in.GetName().GetValue()
		description = in.GetDescription().GetValue()
	)

	success, err := api.pizzaLand.UpdateCategory(ctx, id, name, description)
	if err != nil {
		// TODO: implement no internal error check
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
		success, err = api.pizzaLand.RemoveCategoryById(ctx, v.CategoryId)
	case *pizzalndv1.RemoveCategoryRequest_CategoryName:
		success, err = api.pizzaLand.RemoveCategoryByName(ctx, v.CategoryName)
	case nil:
		return nil, status.Error(codes.InvalidArgument, NoIdentifier)
	default:
		return nil, status.Error(codes.InvalidArgument, UnknownNameOrId)
	}

	if err != nil {
		// TODO: implement no internal error check
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pizzalndv1.RemoveCategoryResponse{Success: success}, nil
}
