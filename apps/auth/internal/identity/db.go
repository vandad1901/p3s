package identity

import (
	"context"
	"time"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	commonpb "github.com/vandad1901/p3s/packages/go/gen/protobuf/commonpb/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (*User) TableName() string {
	return "user_t"
}

func dbCreateUser(tx *gorm.DB, user *User) (*commonpb.IDVersion, error) {
	currentTime := time.Now()

	user.ID = 0
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	err := tx.Create(user).Error
	if err != nil {
		if dbpattern.IsConstraintViolation(err, "user_t_username_key") {
			return nil, ErrDuplicateUsername
		}

		if dbpattern.IsConstraintViolation(err, "user_t_email_key") {
			return nil, ErrDuplicateEmail
		}

		return nil, err
	}

	return &commonpb.IDVersion{
		Id:        user.ID,
		UpdatedAt: timestamppb.New(currentTime)}, nil
}

func dbGetUserByID(_ context.Context, tx *gorm.DB, id int64) (*User, error) {
	var res *User

	err := tx.
		Where("id IN ?", id).
		Select(`
			id,

			username,
			email,
			salt,
			password_hash,
			display_name
		`).
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func dbGetUsersByIDs(_ context.Context, tx *gorm.DB, ids []int64) ([]*UserGetResponse, error) {
	var res []*UserGetResponse

	err := tx.Model(&User{}).
		Where("id IN ?", ids).
		Select(`
			id,
			display_name
		`).
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func dbGetUserByUsername(tx *gorm.DB, username string) (*User, error) {
	user := new(User)

	err := tx.
		Where("username = ?", username).
		Select(`
			id,

			username,
			email,
			salt,
			password_hash,

			created_at,
			updated_at
		`).
		Take(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func dbUpdateUser(_ context.Context, tx *gorm.DB, user *User) error {
	currentTime := time.Now()

	err := tx.Model(&User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"display_name": user.DisplayName,
			"updated_at":   currentTime,
		}).Error
	if err != nil {
		return err
	}

	return nil
}
