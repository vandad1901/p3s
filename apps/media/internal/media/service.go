package media

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/idv"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

type Service struct {
	logger   *slog.Logger
	s3Client *s3.Client
	db       *gorm.DB

	consumer *rabbitmq.Consumer
}

func NewService(logger *slog.Logger, db *gorm.DB, s3Client *s3.Client, consumer *rabbitmq.Consumer) *Service {
	return &Service{
		logger:   logger,
		db:       db,
		s3Client: s3Client,

		consumer: consumer,
	}
}

func (s *Service) takeLease(ctx context.Context, key string) (time.Time, error) {
	db := s.db.WithContext(ctx)

	var (
		res time.Time
		err error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbTakeLease(ctx, tx, key)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return time.Time{}, txErr
	}

	return res, nil
}

func (s *Service) changeStatus(ctx context.Context,
	mediaKey string, queuedAt time.Time, targetStatus MediaIngestStatus,
) (*idv.IDV, error) {
	db := s.db.WithContext(ctx)

	var (
		res *idv.IDV
		err error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		err = dbChangeStatus(ctx, tx, mediaKey, queuedAt, targetStatus)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}
