package dbutil

import (
	"fmt"
	"purpl3shadow/utils/envutil"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenDatabaseConnection() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Toronto",
		envutil.DBHost, envutil.DBUser, envutil.DBPass, envutil.DBName, envutil.DBPort)

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
