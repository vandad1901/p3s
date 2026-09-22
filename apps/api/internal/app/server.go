package app

import (
	"github.com/vandad1901/p3s/apps/api/internal/config"
	postrpc "github.com/vandad1901/p3s/apps/api/internal/post/rpc"
	"github.com/vandad1901/p3s/packages/go/authguard"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func initializeServers(a *App, cfg *config.Config) {
	a.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(
		authguard.GRPCAuthGuardWithExceptions(a.logger, a.parser, a.keyfunc, map[string]struct{}{
			"/api.postpb.v1/Get": {},
		}),
	))
	registerGRPCServers(a, cfg)
}

func registerGRPCServers(a *App, cfg *config.Config) {
	healthpb.RegisterHealthServer(a.grpcServer, health.NewServer())
	postrpc.Register(a.grpcServer, a.PostService)

	if cfg.Environment == envutil.Development {
		reflection.Register(a.grpcServer)
	}
}
