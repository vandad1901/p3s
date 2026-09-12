package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/apps/media/internal/config"
	"github.com/vandad1901/p3s/apps/media/internal/media"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

type App struct {
	logger        *slog.Logger
	s3Client      *s3.Client
	db            *gorm.DB
	rmqConn       *rabbitmq.Conn
	mediaConsumer *rabbitmq.Consumer

	mediaService *media.Service

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

	return a, nil
}

func MustBoot(cfg *config.Config) *App {
	a, err := Boot(cfg, slog.Default())
	if err != nil {
		log.Fatalf("[!] Failed to boot the application: %v", err)
	}

	return a
}

func (a *App) Serve(_ *config.Config) error {
	servers := []func() error{
		func() error {
			a.logger.Info("Starting media consumer loop")

			return a.mediaService.RunLoop()
		},
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

const shutdownTimeoutSecs = 10

func (a *App) Shutdown() {
	a.shutdownOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecs*time.Second)
		defer cancel()

		var shutdownWG sync.WaitGroup

		shutdownWG.Go(func() {
			if a.mediaConsumer != nil {
				a.logger.Info("Shutting down consumer")
				a.mediaConsumer.CloseWithContext(ctx)
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
