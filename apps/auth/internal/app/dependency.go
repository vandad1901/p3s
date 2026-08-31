package app

import (
	"fmt"

	"github.com/MicahParks/jwkset"
	"github.com/vandad1901/p3s/apps/auth/internal/config"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"gorm.io/gorm"
)

func initializeDependencies(cfg *config.Config) (*App, error) {
	db, err := initializeDatabase(cfg)
	if err != nil {
		return nil, err
	}

	signer := token.NewECDSASigner(cfg.JWTConfig.PrivateKey)

	return &App{
		db:     db,
		signer: signer,
		KeySet: jwkset.NewMemoryStorage(),
	}, nil
}

func initializeDatabase(cfg *config.Config) (*gorm.DB, error) {
	db := dbpattern.OpenDatabaseConnection(cfg.DSN)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
