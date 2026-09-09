package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/apps/auth/internal/authn"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/jwks"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	logger *slog.Logger
	db     *gorm.DB
	signer token.Signer
	KeySet *jwkset.MemoryJWKSet

	tokenService    *token.Service
	identityService *identity.Service
	SessionService  *session.Service
	AuthnService    *authn.Service
	JWKSService     *jwks.Service

	grpcServer *grpc.Server
	echo       *echo.Echo

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

	initializeServices(a, cfg.JWTConfig)

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

func (a *App) Serve(cfg *config.Config) error {
	servers := []func() error{
		func() error { return serveGRPC(a, cfg) },
		func() error { return serveHTTP(a, cfg) },
	}

	var runnerWG sync.WaitGroup

	errChan := make(chan error, len(servers))

	for _, srv := range servers {
		runnerWG.Go(func() {
			errChan <- srv()
		})
	}

	firstErr := <-errChan

	a.Shutdown()

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

func serveGRPC(a *App, cfg *config.Config) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(context.Background(), "tcp", cfg.GRPCListenAddress)
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

func serveHTTP(a *App, cfg *config.Config) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(context.Background(), "tcp", cfg.HTTPListenAddress)
	if err != nil {
		return fmt.Errorf("http listen on %s: %w", cfg.HTTPListenAddress, err)
	}

	a.logger.Info("HTTP listening", "address", cfg.HTTPListenAddress)

	a.echo.Listener = lis
	a.echo.HideBanner = true
	a.echo.HidePort = true

	err = a.echo.Start(cfg.HTTPListenAddress)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	a.logger.Info("HTTP server stopped gracefully")

	return nil
}

const shutdownTimeoutSecs = 10

func (a *App) Shutdown() {
	a.shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecs*time.Second)
		defer cancel()

		var shutdownWG sync.WaitGroup

		shutdownWG.Go(func() {
			a.logger.Info("Shutting down gRPC server")

			a.grpcServer.GracefulStop()
		})
		shutdownWG.Go(func() {
			a.logger.Info("Shutting down HTTP server")

			err := a.echo.Shutdown(ctx)
			if err != nil {
				a.logger.Error("Failed to shutdown HTTP server gracefully", "error", err)
			}
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

			err := a.echo.Close()
			if err != nil {
				a.logger.Error("Failed to close HTTP server", "error", err)
			}
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
