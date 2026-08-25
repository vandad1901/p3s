package session

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/vandad1901/p3s/apps/auth/internal/token"
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
		return nil, errMissingMetadata
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
		return nil, fmt.Errorf("issuing jwt: %w", err)
	}

	refreshToken, err := token.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
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

func (s *Service) CheckRefreshTokenTx(ctx context.Context, tx *gorm.DB,
	sessionID, userID int64,
	refreshTokenHash string,
) (bool, error) {
	valid, err := dbCheckRefreshTokenHash(tx, sessionID, userID, refreshTokenHash)
	if err != nil {
		return false, err
	}

	return valid, nil
}
