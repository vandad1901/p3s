package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
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

	UserID int64 `json:"user_id"`
}

func (s *Service) GenerateJWT(userID int64) (string, error) {
	const ttl = time.Minute * 15
	currentTime := time.Now()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  "auth-service",
			Subject: strconv.FormatInt(userID, 10),

			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(ttl)),
			ID:        uuid.New().String(),
		},
		UserID: userID,
	}

	return s.signer.SignedString(claims)
}

func GenerateRefreshToken() (string, error) {
	const refreshTokenLength = 32

	b := make([]byte, refreshTokenLength)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))

	return base64.RawURLEncoding.EncodeToString(sum[:])
}
