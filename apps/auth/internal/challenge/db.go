package challenge

import (
	"context"
	"time"

	"gorm.io/gorm"
)

func dbCreate(_ context.Context, db *gorm.DB, p *Challenge) (int64, error) {
	p.ID = 0

	err := db.Create(p).Error
	if err != nil {
		return 0, err
	}

	return p.ID, nil
}

func dbGetByID(_ context.Context, db *gorm.DB, challengeID int64, challengeType string) (*Challenge, error) {
	res := new(Challenge)

	err := db.
		Where("id = ?", challengeID).
		Where("challenge_type = ?", challengeType).
		Take(res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func dbRevokeChallengeByType(_ context.Context, db *gorm.DB, userID int64, challengeType string) error {
	currentTime := time.Now()

	err := db.
		Model(&Challenge{}).
		Where("user_id = ?", userID).
		Where("challenge_type = ?", challengeType).
		Where("revoked_at IS NULL").
		Update("revoked_at", currentTime).Error
	if err != nil {
		return err
	}

	return nil
}

func dbRevokeByUser(_ context.Context, db *gorm.DB, userID int64) error {
	currentTime := time.Now()

	err := db.
		Model(&Challenge{}).
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Update("revoked_at", currentTime).Error
	if err != nil {
		return err
	}

	return nil
}
