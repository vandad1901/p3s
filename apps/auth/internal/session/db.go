package session

import (
	"gorm.io/gorm"
)

func (_ *Session) TableName() string {
	return "sessions"
}

func dbCreateSession(tx *gorm.DB, session *Session) error {
	if err := tx.Create(session).Error; err != nil {
		return err
	}

	return nil
}

func dbCheckRefreshTokenHash(tx *gorm.DB,
	SessionID, userID int64,
	refreshTokenHash string) (bool, error) {
	q := tx.Table("sessions").
		Where("id = ?", SessionID).
		Where("user_id = ?", userID).
		Where("refresh_token_hash = ?", refreshTokenHash)

	var valid bool

	err := q.Raw("SELECT EXISTS (?)", q).
		Scan(&valid).Error
	if err != nil {
		return false, err
	}

	return valid, nil
}
