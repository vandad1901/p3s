package authguard

import (
	"errors"

	"github.com/vandad1901/p3s/packages/go/apperror"
)

var (
	ErrInvalidAuth = apperror.Unauthenticated("authguard.invalidAuthorization")
)

var (
	errInvalidJWT = errors.New("authguard.invalidJWT")
	errInvalidID  = errors.New("authguard.invalidID")
)
