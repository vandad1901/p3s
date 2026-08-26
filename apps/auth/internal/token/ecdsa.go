package token

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ECDSASigner struct {
	priv *ecdsa.PrivateKey
}

func NewECDSASigner(priv *ecdsa.PrivateKey) *ECDSASigner {
	return &ECDSASigner{
		priv: priv,
	}
}

func (s *ECDSASigner) SignedString(claims *Claims) (string, error) {
	res, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(s.priv)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return res, nil
}
