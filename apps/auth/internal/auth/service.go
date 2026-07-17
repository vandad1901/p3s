package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/vandad1901/p3s/apps/auth/internal/credential"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/packages/go/apperror"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	userpb "github.com/vandad1901/p3s/packages/go/gen/protobuf/userpb/v1"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB

	identityService *identity.Service
	sessionService  *session.Service
}

func NewAuthService(db *gorm.DB,
	identityService *identity.Service, sessionService *session.Service) *Service {
	return &Service{
		db: db,

		identityService: identityService,
		sessionService:  sessionService,
	}
}

func (s *Service) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	db := s.db.WithContext(ctx)

	user, err := identity.MapToUser(req.GetUser())
	if err != nil {
		return nil, err
	}

	err = identity.ValidateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	err = credential.ValidatePassword(req.GetPassword())
	if err != nil {
		return nil, err
	}

	saltBytes, saltStr, err := credential.GenerateSalt()
	if err != nil {
		return nil, err
	}

	passwordHash := credential.HashPasswordArgon2(req.GetPassword(), saltBytes)

	user.Salt = saltStr
	user.PasswordHash = passwordHash

	sessionResponse := new(userpb.AuthSessionResponse)

	err = dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		_, err = identity.CreateUserTx(ctx, tx, user)
		if err != nil {
			return err
		}

		res, err := s.sessionService.CreateSessionForUserTx(ctx, tx, user.ID)
		if err != nil {
			return err
		}

		sessionResponse = mapToSessionResponsePB(res)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userpb.RegisterResponse{
		Session: sessionResponse,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	db := s.db.WithContext(ctx)

	sessionResponse := new(userpb.AuthSessionResponse)

	err := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		user, err := identity.GetUserByUsernameTx(ctx, tx, req.Username)
		if err != nil {
			return err
		}

		salt, err := base64.RawStdEncoding.DecodeString(user.Salt)
		if err != nil {
			return fmt.Errorf("error decoding salt: %s", user.Salt)
		}

		passwordHash := credential.HashPasswordArgon2(req.Password, salt)

		if subtle.ConstantTimeCompare([]byte(user.PasswordHash), []byte(passwordHash)) != 1 {
			return apperror.Unauthenticated("auth.InvalidLogin")
		}

		res, err := s.sessionService.CreateSessionForUserTx(ctx, tx, user.ID)
		if err != nil {
			return err
		}

		sessionResponse = mapToSessionResponsePB(res)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{Session: sessionResponse}, nil
}
