package media

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/idv"
	"gorm.io/gorm"
)

type Service struct {
	s3Client *s3.Client
	db       *gorm.DB
}

func NewService(db *gorm.DB, s3Client *s3.Client) *Service {
	return &Service{
		db:       db,
		s3Client: s3Client,
	}
}

func (s *Service) MediaIngested(ctx context.Context, keys []string) (int64, error) {
	db := s.db.WithContext(ctx)

	var res int64

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		readyCount, err := dbCountReady(tx, keys)
		if err != nil {
			return err
		}

		res = int64(len(keys)) - readyCount

		return nil
	})
	if txErr != nil {
		return 0, txErr
	}

	return res, nil
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
