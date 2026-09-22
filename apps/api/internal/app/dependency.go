package app

import (
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vandad1901/p3s/apps/api/internal/config"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/envutil"
	"github.com/vandad1901/p3s/packages/go/gormslog"
)

func initializeDependencies(a *App, cfg *config.Config) error {
	err := initializeDatabase(a, cfg)
	if err != nil {
		return err
	}

	if cfg.Environment != envutil.Test {
		err = initializeJWT(a, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func initializeDatabase(a *App, cfg *config.Config) error {
	var err error

	a.db, err = dbpattern.OpenDatabaseConnection(cfg.DSN, gormslog.New(a.logger))
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}

	sqlDB, err := a.db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

func initializeJWT(a *App, cfg *config.Config) error {
	var err error

	a.keyfunc, err = keyfunc.NewDefault([]string{cfg.AuthServiceAddress})
	if err != nil {
		return fmt.Errorf("failed to create keyfunc: %w", err)
	}

	a.parser = jwt.NewParser(
		jwt.WithValidMethods([]string{"ES256"}),
		jwt.WithIssuer("auth-service"),
		jwt.WithExpirationRequired(),
	)

	return nil
}
