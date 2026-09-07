package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vandad1901/p3s/apps/media/internal/app"
	"github.com/vandad1901/p3s/apps/media/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	a, err := app.Boot(cfg, logger)
	if err != nil {
		logger.Error("Failed to boot the application",
			"error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErrChan := a.Serve(ctx)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case err := <-runErrChan:
		if err != nil {
			logger.Error("Error occurred while running the application",
				"error", err)
		}

		cancel()
		a.Shutdown()
	case <-sigChan:
		logger.Info("received termination signal")
		cancel()
		a.Shutdown()
	}
}
