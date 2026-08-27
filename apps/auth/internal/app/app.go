package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/labstack/echo/v4"
	"github.com/vandad1901/p3s/apps/auth/internal/authn"
	authnrpc "github.com/vandad1901/p3s/apps/auth/internal/authn/rpc"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/jwks"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const shutdownTimeoutSecs = 10

type App struct {
	db     *gorm.DB
	signer token.Signer
	KeySet *jwkset.MemoryJWKSet

	TokenService    *token.Service
	IdentityService *identity.Service
	SessionService  *session.Service
	AuthnService    *authn.Service
	JWKSService     *jwks.Service

	grpcServer *grpc.Server
	httpServer *echo.Echo
}

func Boot(cfg *config.Config) (*App, error) {
	a, err := initializeResources(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize resources: %w", err)
	}

	initializeV1Services(a, cfg.JWTConfig)

	if cfg.Environment != config.Test {
		initializeServers(a, cfg)
	}

	return a, nil
}

func ForceBoot(cfg *config.Config) *App {
	a, err := Boot(cfg)
	if err != nil {
		log.Fatalf("[!] Failed to boot the application: %v", err)
	}

	return a
}

func initializeResources(cfg *config.Config) (*App, error) {
	db := dbpattern.OpenDatabaseConnection(cfg.DSN)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	signer := token.NewECDSASigner(cfg.JWTConfig.PrivateKey)

	return &App{
		db:     db,
		signer: signer,
		KeySet: jwkset.NewMemoryStorage(),
	}, nil
}

func initializeV1Services(a *App, jwtConfig *config.JWTConfig) {
	a.TokenService = token.NewService(a.signer)
	a.IdentityService = identity.NewService(a.db)
	a.SessionService = session.NewService(a.db, a.TokenService)
	a.AuthnService = authn.NewAuthNService(a.db, a.IdentityService, a.SessionService, a.TokenService)

	a.JWKSService = jwks.NewService(a.KeySet, jwtConfig.PrivateKey, jwtConfig.KeyID)
}

func initializeServers(a *App, cfg *config.Config) {
	a.grpcServer = grpc.NewServer()
	registerGRPCServers(a, a.grpcServer, cfg)

	a.httpServer = echo.New()
	registerHTTPHandlers(a, a.httpServer)
}

func registerGRPCServers(a *App, grpcServer *grpc.Server, cfg *config.Config) {
	authnrpc.Register(grpcServer,
		a.AuthnService, a.IdentityService, a.SessionService)

	if cfg.Environment == config.Development {
		reflection.Register(grpcServer)
	}
}

func registerHTTPHandlers(a *App, httpServer *echo.Echo) {
	httpServer.GET("/jwks.json", a.JWKSService.Handle)
}

const RUNNER_COUNT = 2

func (a *App) Run(cfg *config.Config) chan error {
	errChan := make(chan error, RUNNER_COUNT)

	go func() {
		err := serveGRPC(cfg, a.grpcServer)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		err := serveHTTP(cfg, a.httpServer)
		if err != nil {
			errChan <- err
		}
	}()

	return errChan
}

func serveGRPC(cfg *config.Config, grpcServer *grpc.Server) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(context.Background(), "tcp", cfg.GRPCListenAddress)
	if err != nil {
		return fmt.Errorf("grpc listen on %s: %w", cfg.GRPCListenAddress, err)
	}

	log.Printf("[i] Auth gRPC Listening on %s", cfg.GRPCListenAddress)

	err = grpcServer.Serve(lis)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("serve grpc: %w", err)
	}

	log.Println("[i] gRPC server stopped gracefully")

	return nil
}

func serveHTTP(cfg *config.Config, e *echo.Echo) error {
	lc := net.ListenConfig{}

	lis, err := lc.Listen(context.Background(), "tcp", cfg.HTTPListenAddress)
	if err != nil {
		return fmt.Errorf("http listen on %s: %w", cfg.HTTPListenAddress, err)
	}

	log.Printf("[i] Auth HTTP Listening on %s", cfg.HTTPListenAddress)

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

func (a *App) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecs*time.Second)
	defer cancel()

	var stoppers sync.WaitGroup

	stoppers.Go(func() {
		log.Println("[i] Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
	})
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

		sqlDB, err := a.db.DB()
		if err != nil {
			log.Printf("[!] Failed to get SQL DB: %v", err)
			close(done)

			return
		}

		err = sqlDB.Close()
		if err != nil {
			log.Printf("[!] Failed to close SQL DB: %v", err)
			close(done)

			return
		}

		close(done)
	}()

	select {
	case <-done:
		log.Println("[i] All servers shutdown gracefully completed")

	case <-ctx.Done():
		log.Println("[!] Graceful shutdown timed out, forcing stop")

		a.grpcServer.Stop()

		err := a.httpServer.Close()
		if err != nil {
			log.Printf("[!] Failed to forcefully close HTTP server: %v", err)
		}
	}
}
