package authn

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/vandad1901/p3s/apps/auth/internal/credential"
	"github.com/vandad1901/p3s/apps/auth/internal/identity"
	"github.com/vandad1901/p3s/apps/auth/internal/session"
	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/dbpattern"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/authnpb/v1"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB

	identityService *identity.Service
	sessionService  *session.Service
	tokenService    *token.Service
}

func NewAuthNService(db *gorm.DB,
	identityService *identity.Service, sessionService *session.Service, tokenService *token.Service) *Service {
	return &Service{
		db: db,

		identityService: identityService,
		sessionService:  sessionService,
		tokenService:    tokenService,
	}
}

func (s *Service) Register(ctx context.Context, user *identity.User, password string,
) (*session.SessionResponse, error) {
	db := s.db.WithContext(ctx)

	err := identity.ValidateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("validating user: %w", err)
	}

	err = credential.ValidatePassword(password)
	if err != nil {
		return nil, fmt.Errorf("validating password: %w", err)
	}

	saltBytes, saltStr, err := credential.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}

	passwordHash := credential.HashPasswordArgon2(password, saltBytes)

	user.Salt = saltStr
	user.PasswordHash = passwordHash

	var res *session.SessionResponse

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		_, err = identity.CreateUserTx(ctx, tx, user)
		if err != nil {
			return fmt.Errorf("creating user: %w", err)
		}

		res, err = s.sessionService.CreateSessionForUserTx(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("creating session: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

func (s *Service) Login(ctx context.Context, username string, password string) (*session.SessionResponse, error) {
	db := s.db.WithContext(ctx)

	var res *session.SessionResponse

	txErr := dbpattern.SerializableTx(db, func(tx *gorm.DB) error {
		user, err := identity.GetUserByUsernameTx(ctx, tx, username)
		if err != nil {
			return fmt.Errorf("getting user: %w", err)
		}

		salt, err := base64.RawStdEncoding.DecodeString(user.Salt)
		if err != nil {
			return fmt.Errorf("decoding salt: %w", err)
		}

		passwordHash := credential.HashPasswordArgon2(password, salt)

		if subtle.ConstantTimeCompare([]byte(user.PasswordHash), []byte(passwordHash)) != 1 {
			return errInvalidAuthn
		}

		res, err = s.sessionService.CreateSessionForUserTx(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("creating session: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return res, nil
}

func (s *Service) RefreshJWT(ctx context.Context, in *authnpb.RefreshJWTRequest) (*authnpb.RefreshJWTResponse, error) {
	db := s.db.WithContext(ctx)

	refreshTokenHash := token.HashRefreshToken(in.GetRefreshToken())

	valid, err := s.sessionService.CheckRefreshTokenTx(ctx, db,
		in.GetSessionId(), in.GetUserId(), refreshTokenHash)
	if err != nil || !valid {
		return nil, errInvalidAuthn
	}

	jwt, err := s.tokenService.GenerateJWT(in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("generating jwt: %w", err)
	}

	return &authnpb.RefreshJWTResponse{
		Jwt: jwt,
	}, nil
}
