package app

import (
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	uploadrpc "github.com/vandad1901/p3s/apps/upload/internal/upload/rpc"
	"github.com/vandad1901/p3s/packages/go/authguard"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func initializeServers(a *App, cfg *config.Config) {
	a.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(authguard.GRPCAuthGuard(a.logger, a.parser, a.keyfunc)))
	registerGRPCServers(a, cfg)
}

func registerGRPCServers(a *App, cfg *config.Config) {
	healthpb.RegisterHealthServer(a.grpcServer, health.NewServer())
	uploadrpc.Register(a.grpcServer, a.UploadService)

	if cfg.Environment == envutil.Development {
		reflection.Register(a.grpcServer)
	}
}
