package token

import (
	"crypto/ecdsa"

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
	return jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(s.priv)
}
