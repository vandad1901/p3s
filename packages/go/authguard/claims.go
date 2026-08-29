package authguard

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
}

func (c *Claims) Validate() error {
	if c.ID == "" {
		return errInvalidID
	}

	return nil
}
