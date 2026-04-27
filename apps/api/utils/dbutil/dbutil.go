package dbutil

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDatabaseConnection() {
	dsn := "host=localhost user=postgres password=postgres " +
		"dbname=purpl3shadow port=5432 sslmode=disable TimeZone=America/Toronto"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error connecting to database: %w", err))
	}

	DB = db
}

func Gorm() *gorm.DB {
	return DB.Session(&gorm.Session{NewDB: true})
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
