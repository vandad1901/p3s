package session

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/vandad1901/p3s/apps/auth/internal/token"
	"github.com/vandad1901/p3s/packages/go/apperror"
	"github.com/vandad1901/p3s/packages/go/gen/protobuf/auth/authnpb/v1"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB

	tokenService *token.Service
}

func NewService(db *gorm.DB, tokenService *token.Service) *Service {
	return &Service{
		db: db,

		tokenService: tokenService}
}

const refreshTTL = time.Minute * 60 * 24 * 30

func (s *Service) CreateSessionForUserTx(ctx context.Context, tx *gorm.DB, userID int64) (*SessionResponse, error) {
	currentTime := time.Now()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("session.MissingMetadata")
	}

	var ipAddress net.IP

	xff := md.Get("x-forwarded-for")
	if len(xff) > 0 {
		ipAddress = net.ParseIP(strings.TrimSpace(
			strings.Split(xff[0], ",")[0],
		))
	}

	var userAgent string

	userAgents := md.Get("user-agent")
	if len(userAgents) > 0 {
		userAgent = userAgents[0]
	}

	jwt, err := s.tokenService.GenerateJWT(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshTokenHash := token.HashRefreshToken(refreshToken)

	expiresAt := currentTime.Add(refreshTTL)
	session := &Session{
		UserID: userID,

		RefreshTokenHash: refreshTokenHash,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		Status:           SessionStatusActive,

		IssuedAt:  currentTime,
		ExpiresAt: expiresAt,
	}

	err = dbCreateSession(tx, session)
	if err != nil {
		return nil, err
	}

	return &SessionResponse{
		SessionID:    session.ID,
		JWT:          jwt,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt}, nil
}

func (s *Service) RefreshJWT(ctx context.Context, in *authnpb.RefreshJWTRequest) (*authnpb.RefreshJWTResponse, error) {
	db := s.db.WithContext(ctx)

	refreshTokenHash := token.HashRefreshToken(in.GetRefreshToken())

	valid, err := dbCheckRefreshTokenHash(db, in.GetSessionId(), in.GetUserId(), refreshTokenHash)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, apperror.NotFound("session.InvalidSession")
	}

	jwt, err := s.tokenService.GenerateJWT(in.GetUserId())
	if err != nil {
		return nil, err
	}

	return &authnpb.RefreshJWTResponse{
		Jwt: jwt,
	}, nil
}
