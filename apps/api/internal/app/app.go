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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const shutdownTimeoutSecs = 10

type App struct {
	PostService *post.Service

	db *gorm.DB

	grpcServer *grpc.Server
}

func Boot(cfg *config.Config) (*App, error) {
	var (
		dsn = config.GetDSN()
	)

	db := dbpattern.OpenDatabaseConnection(dsn)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	a := initializeServices(db)

	a.db = db

	if cfg.Environment != config.Test {
		a.grpcServer = grpc.NewServer()
		registerGRPCServers(a, a.grpcServer, cfg)
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

func initializeServices(db *gorm.DB) *App {
	postServiceV1 := post.NewService(db)

	return &App{
		PostService: postServiceV1,
	}
}

func registerGRPCServers(a *App, grpcServer *grpc.Server, cfg *config.Config) {
	postrpc.Register(grpcServer, a.PostService)

	if cfg.Environment == config.Development {
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
