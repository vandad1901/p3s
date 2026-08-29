package app

import (
	"github.com/vandad1901/p3s/apps/api/internal/config"
	postrpc "github.com/vandad1901/p3s/apps/api/internal/post/rpc"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initializeServers(a *App, cfg *config.Config) {
	a.grpcServer = grpc.NewServer()
	registerGRPCServers(a, a.grpcServer, cfg)
}

func registerGRPCServers(a *App, grpcServer *grpc.Server, cfg *config.Config) {
	postrpc.Register(grpcServer, a.PostService)

	if cfg.Environment == envutil.Development {
		reflection.Register(grpcServer)
	}
}
