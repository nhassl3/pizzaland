package app

import (
	"log/slog"
	"time"

	"github.com/nhassl3/pizzaland/internals/app/grpcapp"
	"github.com/nhassl3/pizzaland/internals/app/httpapp"
	"github.com/nhassl3/pizzaland/internals/domain/services/pizzaland"
	"github.com/nhassl3/pizzaland/internals/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
	HTTPServer *httpapp.App
}

func MustLoadApp(
	log *slog.Logger,
	timeout time.Duration,
	gRPCHost string,
	gRPCPort int,
	HTTPHost string,
	HTTPPort int,
	storagePath string,
) *App {
	storage, err := sqlite.NewStorage(timeout, storagePath)
	if err != nil {
		panic(err)
	}

	urlPizzaLandObj := pizzaland.NewPizzaLand(log, storage, storage, storage, storage)

	return &App{
		GRPCServer: grpcapp.NewApp(log, gRPCHost, gRPCPort, urlPizzaLandObj),
		HTTPServer: httpapp.NewApp(log, HTTPHost, HTTPPort, gRPCHost, gRPCPort, timeout),
	}
}
