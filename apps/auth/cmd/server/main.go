package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/apps/auth/internal/auth"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	userrpc "github.com/vandad1901/p3s/apps/auth/internal/identity/rpc"
	"github.com/vandad1901/p3s/apps/auth/internal/jwks"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const shutdownTimeoutSecs = 10

type app struct {
	tokenService    *token.Service
	identityService *identity.Service
	sessionService  *session.Service
	authService     *auth.Service
	jwksService     *jwks.Service
}

func main() {
	var (
		cfg       = config.LoadConfig()
		dsn       = config.GetDSN()
		jwtConfig = config.LoadJWTConfig()
	)

	var (
		grpcServer = grpc.NewServer()
		httpServer = echo.New()
		db         = dbpattern.OpenDatabaseConnection(dsn)
		signer     = token.NewECDSASigner(jwtConfig.PrivateKey)
	)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[!] Failed to Get SQL DB: %v", err)
	}
	defer sqlDB.Close()

	a := initializeServices(jwtConfig, signer, db)

	registerGRPCServers(a, grpcServer, cfg)
	registerHTTPHandlers(a, httpServer)

	manageAppLifecycle(cfg, grpcServer, httpServer)
}

func initializeServices(jwtConfig *config.JWTConfig, signer *token.ECDSASigner, db *gorm.DB) *app {
	tokenServiceV1 := token.NewService(signer)
	identityServiceV1 := identity.NewService(db)
	sessionServiceV1 := session.NewService(db, tokenServiceV1)
	authServiceV1 := auth.NewAuthService(db, identityServiceV1, sessionServiceV1)

	jwksServiceV1 := jwks.NewService(jwtConfig.PrivateKey, "your-key-id")

	return &app{
		tokenService:    tokenServiceV1,
		identityService: identityServiceV1,
		sessionService:  sessionServiceV1,
		authService:     authServiceV1,
		jwksService:     jwksServiceV1,
	}
}

func registerGRPCServers(a *app, grpcServer *grpc.Server, cfg *config.Config) {
	userrpc.Register(grpcServer,
		a.identityService, a.authService, a.sessionService)

	if cfg.Environment == config.Development {
		reflection.Register(grpcServer)
	}
}

func registerHTTPHandlers(a *app, httpServer *echo.Echo) {
	httpServer.GET("/jwks.json", a.jwksService.Handle)
}

func serveGRPC(cfg *config.Config, grpcServer *grpc.Server) {
	lis, err := net.Listen("tcp", cfg.GRPCListenAddress)
	if err != nil {
		log.Fatalf("[!] Failed to listen on gRPC port: %v", err)
	}

	log.Printf("[i] Auth gRPC Listening on %s", cfg.GRPCListenAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[!] Failed to serve gRPC server: %v", err)
	}
}

func serveHTTP(cfg *config.Config, e *echo.Echo) {
	log.Printf("[i] Auth HTTP Listening on %s", cfg.HTTPListenAddress)

	e.HideBanner = true
	e.HidePort = true
	if err := e.Start(cfg.HTTPListenAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[!] Failed to serve HTTP server: %v", err)
	}
}

func manageAppLifecycle(cfg *config.Config, grpcServer *grpc.Server, httpServer *echo.Echo) {
	var wg sync.WaitGroup

	wg.Go(func() {
		serveGRPC(cfg, grpcServer)
	})
	wg.Go(func() {
		serveHTTP(cfg, httpServer)
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
	wg.Go(func() {
		log.Println("[i] Shutting down HTTP server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("[!] Failed to shutdown HTTP server gracefully: %v", err)
		}
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
