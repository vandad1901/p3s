package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/vandad1901/p3s/apps/api/internal/config"
	"github.com/vandad1901/p3s/apps/api/internal/post"
	postrpc "github.com/vandad1901/p3s/apps/api/internal/post/rpc"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB

	PostService *post.Service

	grpcServer *grpc.Server
}

func Boot(cfg *config.Config) (*App, error) {
	a, err := initializeResources(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize resources: %w", err)
	}

	initializeV1Services(a)

	if cfg.Environment != envutil.Test {
		initializeServers(a, cfg)
	}

	return a, nil
}

func MustBoot(cfg *config.Config) *App {
	a, err := Boot(cfg)
	if err != nil {
		log.Fatalf("[!] Failed to boot the application: %v", err)
	}

	return a
}

func initializeResources(cfg *config.Config) (*App, error) {
	db := dbpattern.OpenDatabaseConnection(cfg.DSN)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &App{
		db: db,
	}, nil
}

func initializeV1Services(a *App) {
	a.PostService = post.NewService(a.db)
}

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

const RUNNER_COUNT = 1

func (a *App) Run(cfg *config.Config) chan error {
	errChan := make(chan error, RUNNER_COUNT)

	go func() {
		err := serveGRPC(cfg, a.grpcServer)
		if err != nil {
			errChan <- err
		}
	}()

	return errChan
}

func serveGRPC(cfg *config.Config, grpcServer *grpc.Server) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(context.Background(), "tcp", cfg.GRPCListenAddress)
	if err != nil {
		return fmt.Errorf("grpc listen on %s: %w", cfg.GRPCListenAddress, err)
	}

	log.Printf("[i] API gRPC Listening on %s", cfg.GRPCListenAddress)

	err = grpcServer.Serve(lis)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve grpc: %w", err)
	}

	log.Println("[i] gRPC server stopped gracefully")

	return nil
}

const shutdownTimeoutSecs = 10

func (a *App) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecs*time.Second)
	defer cancel()

	var stoppers sync.WaitGroup

	stoppers.Go(func() {
		log.Println("[i] Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
	})

	done := make(chan struct{})

	go func() {
		stoppers.Wait()

		sqlDB, err := a.db.DB()
		if err != nil {
			log.Fatalf("Failed to get SQL DB: %v", err)
			close(done)

			return
		}

		err = sqlDB.Close()
		if err != nil {
			log.Fatalf("Failed to close SQL DB: %v", err)
			close(done)

			return
		}

		close(done)
	}()

	select {
	case <-done:
		log.Println("[i] All servers shutdown gracefully completed")
	case <-ctx.Done():
		log.Println("[!] Graceful shutdown timed out, forcing stop")
		a.grpcServer.Stop()
	}
}
