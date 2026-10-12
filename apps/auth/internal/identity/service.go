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

func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
	db := s.db.WithContext(ctx)

	var (
		res *User
		err error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbGetUserByID(ctx, tx, id)
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

func (s *Service) GetUsers(ctx context.Context, ids []int64) ([]*UserGetResponse, error) {
	db := s.db.WithContext(ctx)

	var (
		res []*UserGetResponse
		err error
	)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		res, err = dbGetUsersByIDs(ctx, tx, ids)
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

func GetUserByUsernameTx(ctx context.Context, tx *gorm.DB, username string) (*User, error) {
	res, err := dbGetUserByUsername(tx, username)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) UpdateUser(ctx context.Context, user *User) error {
	db := s.db.WithContext(ctx)

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		err := dbUpdateUser(ctx, tx, user)
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
