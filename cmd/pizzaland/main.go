package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	log.Info("Starting pizzaland service", 
		slog.Int("grpc_port", cfg.GRPC.Port),
		slog.Int("http_port", cfg.HTTP.Port))

	application := app.MustLoadApp(
		log,
		cfg.GRPC.Timeout,
		cfg.GRPC.Host,
		cfg.GRPC.Port,
		cfg.HTTP.Host,
		cfg.HTTP.Port,
		cfg.StoragePath,
	)

	// Start gRPC server first in a goroutine
	go func() {
		application.GRPCServer.MustStart()
	}()
	
	// Wait a bit for gRPC server to start listening
	// This gives gRPC server time to bind to the port
	time.Sleep(500 * time.Millisecond)
	
	// Then start HTTP server
	go func() {
		application.HTTPServer.MustStart()
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	log.Info("Pizzaland server stopped", slog.String("signal", (<-sig).String()))
}
