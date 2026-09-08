package dbpattern

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const PostgresDSNFormat = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC"

func OpenDatabaseConnection(dsn string, logger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger})
	if err != nil {
		return nil, fmt.Errorf("gorm opening: %w", err)
	}

	return db, nil
}

func SerializableTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin(&sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	err := tx.Error
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err = fn(tx)
	if err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}
