package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vandad1901/p3s/apps/auth/internal/app"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	a, err := app.Boot(cfg)
	if err != nil {
		log.Fatalf("[!] Failed to boot the application: %v", err)
	}

	runErrChan := a.Run(cfg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-runErrChan:
		log.Printf("[!] Error occurred while running the application: %v", err)
		a.Close()
	case sig := <-sigChan:
		log.Printf("[i] Received signal %v, shutting down gracefully...", sig)
		a.Close()
	}
}
