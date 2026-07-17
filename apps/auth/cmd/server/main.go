package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vandad1901/p3s/apps/auth/internal/auth"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	userrpc "github.com/vandad1901/p3s/apps/auth/internal/identity/rpc"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const shutdownTimeoutSecs = 10

func main() {
	var (
		cfg       = config.LoadConfig()
		dsn       = config.GetDSN()
		jwtConfig = config.LoadJWTConfig()
	)

	var (
		grpcServer = grpc.NewServer()
		db         = dbpattern.OpenDatabaseConnection(dsn)
		signer     = token.NewECDSASigner(jwtConfig.PrivateKey)
	)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[!] Failed to Get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	registerGRPCServers(signer, db, grpcServer, cfg)

	handleGRPCServer(cfg, grpcServer)
}

func registerGRPCServers(signer *token.ECDSASigner, db *gorm.DB, grpcServer *grpc.Server, cfg *config.Config) {
	tokenServiceV1 := token.NewService(signer)
	identityServiceV1 := identity.NewService(db)
	sessionServiceV1 := session.NewService(db,
		tokenServiceV1)
	authServiceV1 := auth.NewAuthService(db,
		identityServiceV1, sessionServiceV1)

	userrpc.RegisterUserRPCServer(grpcServer,
		identityServiceV1, authServiceV1, sessionServiceV1)

	if cfg.Environment == config.Development {
		reflection.Register(grpcServer)
	}
}

func handleGRPCServer(cfg *config.Config, grpcServer *grpc.Server) {
	lis, err := net.Listen("tcp", cfg.GRPCListenAddress)
	if err != nil {
		log.Fatalf("[!] Failed to listen: %v", err)
	}

	go func() {
		log.Printf("[i] Auth Listening on %s", cfg.GRPCListenAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("[!] Failed to serve gRPC server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("[i] Received signal %v, shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecs*time.Second)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("[i] Server shutdown gracefully completed")
	case <-ctx.Done():
		log.Println("[!] Graceful shutdown timed out, forcing stop")
		grpcServer.Stop()
	}
}
