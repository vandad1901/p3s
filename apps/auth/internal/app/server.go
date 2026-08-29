package app

import (
	"github.com/labstack/echo/v4"
	authnrpc "github.com/vandad1901/p3s/apps/auth/internal/authn/rpc"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	jwkshttp "github.com/vandad1901/p3s/apps/auth/internal/jwks/http"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initializeServers(a *App, cfg *config.Config) {
	a.grpcServer = grpc.NewServer()
	// TODO: Add gRPC middleware here
	registerGRPCServers(a, a.grpcServer, cfg)

	a.httpServer = echo.New()
	registerHTTPHandlers(a, a.httpServer)
}

func registerGRPCServers(a *App, grpcServer *grpc.Server, cfg *config.Config) {
	authnrpc.Register(grpcServer,
		a.AuthnService, a.identityService, a.SessionService)

	if cfg.Environment == envutil.Development {
		reflection.Register(grpcServer)
	}
}

func registerHTTPHandlers(a *App, e *echo.Echo) {
	jwkshttp.Register(e, a.JWKSService)
}
