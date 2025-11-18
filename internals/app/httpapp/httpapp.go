package httpapp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/nhassl3/pizzaland/internals/http/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const opStart = "httpapp.MustStart"

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	host       string
	port       int
}

func NewApp(
	log *slog.Logger,
	host string,
	httpPort int,
	grpcHost string,
	grpcPort int,
	timeout time.Duration,
) *App {
	// Create gRPC connection (lazy connection - will connect when needed)
	grpcAddr := fmt.Sprintf("%s:%d", grpcHost, grpcPort)
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC client: %w", err))
	}

	// Create handlers
	h := handlers.NewHandler(log, conn, timeout)

	mux := http.NewServeMux()
	
	// Pizza endpoints
	mux.HandleFunc("GET /api/pizzas", h.ListPizzas)
	mux.HandleFunc("GET /api/pizzas/{id}", h.GetPizza)
	mux.HandleFunc("POST /api/pizzas", h.SavePizza)
	mux.HandleFunc("PUT /api/pizzas/{id}", h.UpdatePizza)
	mux.HandleFunc("DELETE /api/pizzas/{id}", h.RemovePizza)
	
	// Category endpoints
	mux.HandleFunc("GET /api/categories", h.ListCategories)
	mux.HandleFunc("GET /api/categories/{id}", h.GetCategory)
	mux.HandleFunc("GET /api/categories/{id}/pizzas", h.GetCategoryPizzas)
	
	// CORS middleware
	handler := corsMiddleware(mux)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, httpPort),
		Handler:      handler,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	return &App{
		log:        log,
		httpServer: httpServer,
		host:       host,
		port:       httpPort,
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (app *App) MustStart() {
	log := app.log.With(slog.String("op", opStart), slog.Int("port", app.port))

	log.Info("HTTP server starting", slog.String("address", app.httpServer.Addr))

	if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("%s: %w", opStart, err))
	}
}

func (app *App) Stop(ctx context.Context) error {
	return app.httpServer.Shutdown(ctx)
}

