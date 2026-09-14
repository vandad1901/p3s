package app

import (
	"github.com/labstack/echo/v4"
	authnrpc "github.com/vandad1901/p3s/apps/auth/internal/authn/rpc"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	jwkshttp "github.com/vandad1901/p3s/apps/auth/internal/jwks/http"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func initializeServers(a *App, cfg *config.Config) {
	initializeGRPC(a, cfg)
	initializeHTTP(a)
}

func initializeGRPC(a *App, cfg *config.Config) {
	a.grpcServer = grpc.NewServer()
	registerGRPCServers(a, a.grpcServer, cfg)
}

func registerGRPCServers(a *App, grpcServer *grpc.Server, cfg *config.Config) {
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())
	authnrpc.Register(grpcServer, a.AuthnService)

	if cfg.Environment == envutil.Development {
		reflection.Register(grpcServer)
	}
}

func initializeHTTP(a *App) {
	a.echo = echo.New()
	registerHTTPHandlers(a, a.echo)
}

func registerHTTPHandlers(a *App, e *echo.Echo) {
	jwkshttp.Register(e, a.JWKSService)
}
