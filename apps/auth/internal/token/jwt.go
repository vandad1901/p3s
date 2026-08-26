package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Signer interface {
	SignedString(claims *Claims) (string, error)
}

type Service struct {
	signer Signer
}

func NewService(signer Signer) *Service {
	return &Service{
		signer: signer,
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (s *Service) GenerateJWT(userID int64) (string, error) {
	const ttl = time.Minute * 15

	currentTime := time.Now()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.FormatInt(userID, 10),

			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(ttl)),
			ID:        uuid.New().String(),
		},
	}

	jwt, err := s.signer.SignedString(claims)
	if err != nil {
		return "", fmt.Errorf("signing jwt: %w", err)
	}

	return jwt, nil
}

func GenerateRefreshToken() (string, error) {
	const refreshTokenLength = 32

	b := make([]byte, refreshTokenLength)

	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("generating random bits: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return base64.RawURLEncoding.EncodeToString(sum[:])
}
