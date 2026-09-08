package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
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
	logger *slog.Logger
	db     *gorm.DB

	PostService *post.Service

	grpcServer *grpc.Server

	shutdownOnce sync.Once
}

func Boot(cfg *config.Config, logger *slog.Logger) (*App, error) {
	a := &App{
		logger: logger,
	}

	err := initializeDependencies(a, cfg)
	if err != nil {
		a.closeConnections()

		return nil, fmt.Errorf("initialize dependencies: %w", err)
	}

	initializeServices(a)

	if cfg.Environment != envutil.Test {
		initializeServers(a, cfg)
	}

	return a, nil
}

func MustBoot(cfg *config.Config) *App {
	a, err := Boot(cfg, slog.Default())
	if err != nil {
		log.Fatalf("[!] Failed to boot the application: %v", err)
	}

	return a
}

func (a *App) Serve(ctx context.Context, cfg *config.Config) error {
	servers := []func() error{
		func() error { return serveGRPC(ctx, a, cfg) },
	}

	var runnerWG sync.WaitGroup

	errChan := make(chan error, len(servers))

	for _, srv := range servers {
		runnerWG.Go(func() {
			errChan <- srv()
		})
	}

	firstErr := <-errChan

	a.Shutdown(ctx)

	runnerWG.Wait()
	close(errChan)

	if firstErr != nil {
		return firstErr
	}

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func serveGRPC(ctx context.Context, a *App, cfg *config.Config) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(ctx, "tcp", cfg.GRPCListenAddress)
	if err != nil {
		return fmt.Errorf("grpc listen on %s: %w", cfg.GRPCListenAddress, err)
	}

	a.logger.Info("gRPC listening", "address", cfg.GRPCListenAddress)

	err = a.grpcServer.Serve(lis)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve grpc: %w", err)
	}

	a.logger.Info("gRPC server stopped gracefully")

	return nil
}

const shutdownTimeoutSecs = 10

func (a *App) Shutdown(ctx context.Context) {
	a.shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeoutSecs*time.Second)
		defer cancel()

		var shutdownWG sync.WaitGroup

		shutdownWG.Go(func() {
			a.logger.Info("Shutting down gRPC server")

			a.grpcServer.GracefulStop()
		})

		done := make(chan struct{})

		go func() {
			shutdownWG.Wait()

			close(done)
		}()

		select {
		case <-done:
			a.logger.Info("Graceful shutdown; Releasing resources")
			a.closeConnections()

		case <-ctx.Done():
			a.logger.Error("Graceful shutdown timed out; Giving up")

			a.grpcServer.Stop()
		}
	})
}

func (a *App) closeConnections() {
	if a.db != nil {
		a.logger.Info("closing database")

		sqlDB, err := a.db.DB()
		if err != nil {
			a.logger.Error("failed to get SQL database during shutdown", "error", err)
		} else {
			err = sqlDB.Close()
			if err != nil {
				a.logger.Error("failed to close database", "error", err)
			}
		}
	}
}
