package mediaupload

import (
	"context"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/idv"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateMedia(ctx context.Context, key string) (*idv.IDV, error) {
	db := s.db.WithContext(ctx)

	var (
		res *idv.IDV
		err error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbCreateWithKey(ctx, tx, key)
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

func (s *Service) ChangeStatus(ctx context.Context, req *idv.IDV, targetStatus MediaUploadStatus) (*idv.IDV, error) {
	db := s.db.WithContext(ctx)

	var (
		res *idv.IDV
		err error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbChangeStatus(ctx, tx, req, targetStatus)
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

func (s *Service) Cleanup(ctx context.Context) error {
	db := s.db.WithContext(ctx)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		err := dbDeleteUnused(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}
