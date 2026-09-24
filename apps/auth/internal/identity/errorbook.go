package identity

import (
	"errors"

	"github.com/vandad1901/p3s/packages/go/apperror"
)

var (
	ErrDuplicateUsername = apperror.Conflict("identity.DuplicateUsername")
	ErrDuplicateEmail    = apperror.Conflict("identity.DuplicateEmail")
)

var (
	errEmptyUsername   = errors.New("identity.validation.emptyUsername")
	errInvalidUsername = errors.New("identity.validation.invalidUsername")
	errEmptyEmail      = errors.New("identity.validation.emptyEmail")
	errInvalidEmail    = errors.New("identity.validation.invalidEmail")
)
