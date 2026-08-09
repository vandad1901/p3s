package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vandad1901/p3s/apps/api/internal/config"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const shutdownTimeoutSecs = 10

type app struct {
}

func main() {
	var (
		cfg = config.LoadConfig()
		dsn = config.GetDSN()
	)

	var (
		grpcServer = grpc.NewServer()
		db         = dbpattern.OpenDatabaseConnection(dsn)
	)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[!] Failed to Get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	a := initializeServices(db)

	registerGRPCServers(a, grpcServer, cfg)

	manageAppLifecycle(cfg, grpcServer)
}

func initializeServices(db *gorm.DB) *app {
	return &app{}
}

func registerGRPCServers(a *app, grpcServer *grpc.Server, cfg *config.Config) {
	if cfg.Environment == config.Development {
		reflection.Register(grpcServer)
	}
}

func serveGRPC(cfg *config.Config, grpcServer *grpc.Server) {
	lis, err := net.Listen("tcp", cfg.GRPCListenAddress)
	if err != nil {
		log.Fatalf("[!] Failed to listen on gRPC port: %v", err)
	}

	log.Printf("[i] API gRPC Listening on %s", cfg.GRPCListenAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[!] Failed to serve gRPC server: %v", err)
	}
}

func manageAppLifecycle(cfg *config.Config, grpcServer *grpc.Server) {
	var wg sync.WaitGroup

	wg.Go(func() {
		serveGRPC(cfg, grpcServer)
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("[i] Received signal %v, shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecs*time.Second)
	defer cancel()

	wg.Go(func() {
		log.Println("[i] Shutting down gRPC server...")
		grpcServer.GracefulStop()
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[i] All servers shutdown gracefully completed")
	case <-ctx.Done():
		log.Println("[!] Graceful shutdown timed out, forcing stop")
		grpcServer.Stop()
	}
}
