package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nhassl3/pizzaland/internals/app"
	"github.com/nhassl3/pizzaland/internals/config"
	"github.com/nhassl3/pizzaland/internals/lib/logger"
)

var (
	cfg *config.Config
	log *slog.Logger
)

func init() {
	cfg = config.MustLoad()

	log = logger.MustLoad(cfg.EnvLevel)
	slog.SetDefault(log)
}

func main() {
	log.Info("Starting pizzaland service", slog.Int("port", cfg.GRPC.Port))

	application := app.MustLoadApp(log, cfg.GRPC.Port, cfg.StoragePath)

	go application.GRPCServer.MustStart()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	log.Info("Pizzaland server stopped", slog.String("signal", (<-sig).String()))
}
