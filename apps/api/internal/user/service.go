package user

import (
	"context"
	"purpl3shadow/gen/commonpb"
	"purpl3shadow/gen/userpb"
	"purpl3shadow/utils/dbutil"

	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*userpb.User, error)
	GetUserByUsername(ctx context.Context, username string) (*userpb.User, error)
	CreateUser(ctx context.Context, user *userpb.User) (*commonpb.IDVersion, error)
	UpdateUser(ctx context.Context, user *userpb.User) (*commonpb.IDVersion, error)
	DeleteUser(ctx context.Context, id int64) error
}

type UserServiceImpl struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{db: db}
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id int64) (*userpb.User, error) {
	var user User

	err := dbutil.SerializableTx(s.db, func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userpb.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*userpb.User, error) {
	var user User

	err := dbutil.SerializableTx(s.db, func(tx *gorm.DB) error {
		if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userpb.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, in *userpb.User) (*commonpb.IDVersion, error) {
	var (
		res *commonpb.IDVersion
	)

	user, err := mapToUser(in)
	if err != nil {
		return nil, err
	}

	err = dbutil.SerializableTx(s.db, func(tx *gorm.DB) error {
		res, err = dbCreateUser(tx, user)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, in *userpb.User) (*commonpb.IDVersion, error) {
	var (
		res *commonpb.IDVersion
	)

	user, err := mapToUser(in)
	if err != nil {
		return nil, err
	}

	err = dbutil.SerializableTx(s.db, func(tx *gorm.DB) error {
		res, err = dbUpdateUser(tx, user)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	return s.db.Delete(&User{}, id).Error
}
