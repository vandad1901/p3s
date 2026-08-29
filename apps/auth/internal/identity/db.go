package identity

import (
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
