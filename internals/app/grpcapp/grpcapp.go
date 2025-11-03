package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/nhassl3/pizzaland/internals/domain/services/pizzaland"
	pizzaLandGRPC "github.com/nhassl3/pizzaland/internals/grpc/pizzaland"
	"google.golang.org/grpc"
)

const opStart = "grpcapp.MustStart"

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	host       string
	port       int
}

func NewApp(log *slog.Logger,
	host string,
	gRPCPort int,
	pizzaLandObj *pizzaland.DomainPizzaLand,
) *App {
	gRPCServer := grpc.NewServer()

	pizzaLandGRPC.Register(gRPCServer, pizzaLandObj)

	return &App{
		gRPCServer: gRPCServer,
		port:       gRPCPort,
		host:       host,
		log:        log,
	}
}

func (app *App) MustStart() {
	log := app.log.With(slog.String("op", opStart), slog.Int("port", app.port))

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", app.host, app.port))
	if err != nil {
		panic(fmt.Errorf("%s: %w", opStart, err))
	}

	log.Info("Server started", slog.String("address", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil {
		panic(fmt.Errorf("%s: %w", opStart, err))
	}
}

func (app *App) Stop() {
	app.gRPCServer.GracefulStop()
}
