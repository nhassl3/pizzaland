package app

import (
	"log/slog"
	"time"

	"github.com/nhassl3/pizzaland/internals/app/grpcapp"
	"github.com/nhassl3/pizzaland/internals/domain/services/pizzaland"
	"github.com/nhassl3/pizzaland/internals/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func MustLoadApp(
	log *slog.Logger,
	timeout time.Duration,
	host string,
	gRPCPort int,
	storagePath string,
) *App {
	storage, err := sqlite.NewStorage(timeout, storagePath)
	if err != nil {
		panic(err)
	}

	urlPizzaLandObj := pizzaland.NewPizzaLand(log, storage, storage, storage, storage)

	return &App{
		GRPCServer: grpcapp.NewApp(log, host, gRPCPort, urlPizzaLandObj),
	}
}
