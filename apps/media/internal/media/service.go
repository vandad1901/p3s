package media

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"gorm.io/gorm"
)

type Service struct {
	s3ExternalEndpoint string
	db                 *gorm.DB
	s3Client           *s3.Client
}

func NewService(s3ExternalEndpoint string, db *gorm.DB, s3Client *s3.Client) *Service {
	return &Service{
		s3ExternalEndpoint: s3ExternalEndpoint,
		db:                 db,
		s3Client:           s3Client,
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

func (s *Service) GetMediaInfos(ctx context.Context, keys []string) ([]MediaInfo, error) {
	db := s.db.WithContext(ctx)

	var (
		res []MediaInfo
		err error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbGetMediaInfos(ctx, tx, keys)
		if err != nil {
			return err
		}

		for _, media := range res {
			for i := range media.Derivatives {
				media.Derivatives[i].URL = fmt.Sprintf("%s/%s/%s.%d.%s",
					s.s3ExternalEndpoint,
					s3Bucket,
					media.MediaKey, media.Derivatives[i].Width, media.Derivatives[i].Ext)
			}
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
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

func (s *Service) finalizeIngest(ctx context.Context,
	media *Media, derivatives []derivativeAsset,
) error {
	db := s.db.WithContext(ctx)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		mediaID, err := dbFinalizeIngest(ctx, tx, media)
		if err != nil {
			return err
		}

		err = dbCreateDerivatives(ctx, tx, mediaID, media.MediaKey, derivatives)
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
