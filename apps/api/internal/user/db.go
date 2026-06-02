package user

import (
	"purpl3shadow/gen/commonpb"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func dbCreateUser(tx *gorm.DB, user *User) (*commonpb.IDVersion, error) {
	currentTime := time.Now()

	user.ID = 0
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	if err := tx.Create(&user).Error; err != nil {
		return nil, err
	}

	return &commonpb.IDVersion{
		Id:        user.ID,
		UpdatedAt: timestamppb.New(currentTime)}, nil
}

func dbUpdateUser(tx *gorm.DB, user *User) (*commonpb.IDVersion, error) {
	currentTime := time.Now()

	user.UpdatedAt = currentTime

	err := tx.Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"username":      user.Username,
			"email":         user.Email,
			"salt":          user.Salt,
			"password_hash": user.PasswordHash,
			"updated_at":    user.UpdatedAt,
		}).Error
	if err != nil {
		return nil, err
	}

	return &commonpb.IDVersion{
		Id:        user.ID,
		UpdatedAt: timestamppb.New(currentTime)}, nil
}
