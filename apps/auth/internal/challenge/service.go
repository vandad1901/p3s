package challenge

import (
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"time"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
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

func (s *Service) Create(ctx context.Context, challenge *Challenge, secret []byte) (int64, error) {
	db := s.db.WithContext(ctx)

	err := validateOnCreate(ctx, db, challenge)
	if err != nil {
		return 0, err
	}

	challenge.SecretHash = sha512.Sum512(secret)
	challenge.ConsumedAt = time.Time{}

	var res int64

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbCreate(ctx, tx, challenge)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return 0, txErr
	}

	return res, nil
}

func (s *Service) Consume(ctx context.Context, challengeID int64, challengeType string, secret []byte) error {
	db := s.db.WithContext(ctx)

	secretHash := sha512.Sum512(secret)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		currentTime := time.Now()

		challenge, err := dbGetByID(ctx, tx, challengeID, challengeType)
		if err != nil {
			return err
		}

		if challenge.ExpiresAt.After(currentTime) {
			return errExpiredChallenge
		}

		if !challenge.ConsumedAt.IsZero() {
			return errAlreadyConsumed
		}

		if ok := subtle.ConstantTimeCompare(secretHash[:], challenge.SecretHash[:]); ok == 0 {
			return errBadHash
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}

func (s *Service) RevokeChallenge(ctx context.Context, userID int64, challengeType string) error {
	db := s.db.WithContext(ctx)

	err := dbRevokeChallengeByType(ctx, db, userID, challengeType)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RevokeAllUserChallenges(ctx context.Context, userID int64) error {
	db := s.db.WithContext(ctx)

	err := dbRevokeByUser(ctx, db, userID)
	if err != nil {
		return err
	}

	return nil
}
