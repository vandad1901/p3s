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

	"github.com/MicahParks/keyfunc/v3"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/apps/upload/internal/config"
	"github.com/vandad1901/p3s/apps/upload/internal/upload"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

type App struct {
	logger    *slog.Logger
	s3Client  *s3.Client
	db        *gorm.DB
	keyfunc   keyfunc.Keyfunc
	parser    *jwt.Parser
	rmqConn   *rabbitmq.Conn
	publisher *rabbitmq.Publisher

	uploadService *upload.Service

	echo *echo.Echo

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

	initializeServices(a, cfg)

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
		func() error { return serveHTTP(ctx, a, cfg) },
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

func serveHTTP(ctx context.Context, a *App, cfg *config.Config) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(ctx, "tcp", cfg.HTTPListenAddress)
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

func (a *App) Shutdown(ctx context.Context) {
	a.shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeoutSecs*time.Second)
		defer cancel()

		var shutdownWG sync.WaitGroup

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

			err := a.echo.Close()
			if err != nil {
				a.logger.Error("Failed to close HTTP server", "error", err)
			}
		}
	})
}

func (a *App) closeConnections() {
	if a.rmqConn != nil {
		a.logger.Info("closing RabbitMQ connection")

		err := a.rmqConn.Close()
		if err != nil {
			a.logger.Error("failed to stop rabbitMQ connection", "error", err)
		}
	}

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
