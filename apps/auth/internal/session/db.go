package session

import (
	"time"

	"gorm.io/gorm"
)

func (*Session) TableName() string {
	return "session"
}

func dbCreateSession(tx *gorm.DB, session *Session) error {
	err := tx.Create(session).Error
	if err != nil {
		return err
	}

	return nil
}

func dbCheckRefreshTokenHash(tx *gorm.DB,
	sessionID, userID int64,
	refreshTokenHash string) (bool, error) {
	q := tx.Model(Session{}).
		Where("id = ?", sessionID).
		Where("user_id = ?", userID).
		Where("refresh_token_hash = ?", refreshTokenHash).
		Where("status = ?", SessionStatusActive).
		Where("expires_at > ?", time.Now())

	var valid bool

	err := tx.Raw("SELECT EXISTS (?)", q).
		Scan(&valid).Error
	if err != nil {
		return false, err
	}

	return valid, nil
}
