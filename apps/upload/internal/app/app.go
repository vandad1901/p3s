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
)

type App struct {
	logger   *slog.Logger
	s3Client *s3.Client
	keyfunc  keyfunc.Keyfunc
	parser   *jwt.Parser

	uploadService *upload.Service

	httpServer *echo.Echo
}

func Boot(cfg *config.Config) (*App, error) {
	a, err := initializeResources(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize resources: %w", err)
	}

	initializeServices(a, cfg)

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
		err := serveHTTP(cfg, a.httpServer)
		if err != nil {
			errChan <- err
		}
	}()

	return errChan
}

func serveHTTP(cfg *config.Config, e *echo.Echo) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(context.Background(), "tcp", cfg.HTTPListenAddress)
	if err != nil {
		return fmt.Errorf("http listen on %s: %w", cfg.HTTPListenAddress, err)
	}

	log.Printf("[i] Upload HTTP Listening on %s", cfg.HTTPListenAddress)

	e.Listener = lis
	e.HideBanner = true
	e.HidePort = true

	err = e.Start(cfg.HTTPListenAddress)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	log.Println("[i] HTTP server stopped gracefully")

	return nil
}

const shutdownTimeoutSecs = 10

func (a *App) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecs*time.Second)
	defer cancel()

	var stoppers sync.WaitGroup

	stoppers.Go(func() {
		log.Println("[i] Shutting down HTTP server...")

		err := a.httpServer.Shutdown(ctx)
		if err != nil {
			log.Printf("[!] Failed to shutdown HTTP server gracefully: %v", err)
		}
	})

	done := make(chan struct{})

	go func() {
		stoppers.Wait()

		close(done)
	}()

	select {
	case <-done:
		log.Println("[i] All servers shutdown gracefully completed")

	case <-ctx.Done():
		log.Println("[!] Graceful shutdown timed out, forcing stop")

		err := a.httpServer.Close()
		if err != nil {
			log.Printf("[!] Failed to forcefully close HTTP server: %v", err)
		}
	}
}
