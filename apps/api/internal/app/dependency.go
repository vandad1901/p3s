package app

import (
	"fmt"

	"github.com/vandad1901/p3s/apps/api/internal/config"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
)

func initializeDependency(cfg *config.Config) (*App, error) {
	db := dbpattern.OpenDatabaseConnection(cfg.DSN)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &App{
		db: db,
	}, nil
}
