package app

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
	"github.com/vandad1901/p3s/apps/api/internal/post"
	postrpc "github.com/vandad1901/p3s/apps/api/internal/post/rpc"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const shutdownTimeoutSecs = 10

type App struct {
	PostService *post.Service

	db *gorm.DB

	grpcServer *grpc.Server
}

func Boot() *App {
	var (
		cfg = config.LoadConfig()
		dsn = config.GetDSN()
	)

	db := dbpattern.OpenDatabaseConnection(dsn)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[!] Failed to get SQL DB: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("[!] Failed to ping database: %v", err)
	}

	a := initializeServices(db)

	a.db = db

	if cfg.Environment != config.Test {
		a.grpcServer = grpc.NewServer()
		registerGRPCServers(a, a.grpcServer, cfg)
	}

	return a
}

func (a *App) Run() {
	if a.grpcServer == nil {
		return
	}

	cfg := config.LoadConfig()

	var wg sync.WaitGroup

	wg.Go(func() {
		serveGRPC(cfg, a.grpcServer)
	})

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	sig := <-sigChan
	log.Printf("[i] Received signal %v, shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeoutSecs*time.Second,
	)
	defer cancel()

	wg.Go(func() {
		log.Println("[i] Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
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
		a.grpcServer.Stop()
	}

	a.Close()
}

func (a *App) Close() {
	sqlDB, err := a.db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB: %v", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("Failed to close SQL DB: %v", err)
	}
}

func initializeServices(db *gorm.DB) *App {
	postServiceV1 := post.NewService(db)

	return &App{
		PostService: postServiceV1,
	}
}

func registerGRPCServers(a *App, grpcServer *grpc.Server, cfg *config.Config) {
	postrpc.Register(grpcServer, a.PostService)

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

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("[!] Failed to serve gRPC server: %v", err)
	}
}
