package dbpattern

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const PostgresDSNFormat = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Toronto"

func OpenDatabaseConnection(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error connecting to database: %w", err))
	}

	return db
}

func SerializableTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}
