package identity

import (
	"context"

	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/commonpb/v1"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateUser(ctx context.Context, user *User) (*commonpb.IDVersion, error) {
	db := s.db.WithContext(ctx)

	err := ValidateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	var idv *commonpb.IDVersion

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		idv, err = CreateUserTx(ctx, tx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return idv, nil
}

func CreateUserTx(ctx context.Context, tx *gorm.DB, user *User) (*commonpb.IDVersion, error) {
	idv, err := dbCreateUser(tx, user)
	if err != nil {
		return nil, err
	}

	return idv, nil
}

func GetUserByUsernameTx(ctx context.Context, tx *gorm.DB, username string) (*User, error) {
	res, err := dbGetUserByUsername(tx, username)
	if err != nil {
		return nil, err
	}

	return res, nil
}
