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
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	db *gorm.DB

	PostService *post.Service

	grpcServer *grpc.Server
}

func Boot(cfg *config.Config) (*App, error) {
	a, err := initializeDependency(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize resources: %w", err)
	}

	initializeServices(a)

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
