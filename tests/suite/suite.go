package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	pizzalndv1 "github.com/nhassl3/pizzaland/api/generated/go/pizzaland"
	"github.com/nhassl3/pizzaland/internals/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gRPCHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg         *config.Config
	PizzaClient pizzalndv1.PizzaLandClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByString("../config/local_tests.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.NewClient(
		net.JoinHostPort(gRPCHost, strconv.Itoa(cfg.GRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed %v", err)
	}

	return ctx, &Suite{
		t,
		cfg,
		pizzalndv1.NewPizzaLandClient(cc),
	}
}
