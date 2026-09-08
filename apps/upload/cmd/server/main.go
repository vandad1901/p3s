package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vandad1901/p3s/apps/upload/internal/app"
	"github.com/vandad1901/p3s/apps/upload/internal/config"
)

func main() {
	ok := run()
	if !ok {
		os.Exit(1)
	}
}

func run() bool {
	var (
		cfg    = config.LoadConfig()
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	)

	a, err := app.Boot(cfg, logger)
	if err != nil {
		logger.Error("Failed to boot the application", "error", err)

		return false
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	go func() {
		sig := <-sigChan
		logger.Info("received termination signal", "signal", sig.String())

		cancel()
		a.Shutdown(ctx)
	}()

	err = a.Serve(ctx, cfg)
	if err != nil {
		logger.Error("failed to serve", "error", err)

		return false
	}

	return true
}
