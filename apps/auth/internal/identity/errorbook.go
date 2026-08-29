package identity

import "errors"

var (
	ErrDuplicateUsername = errors.New("identity.DuplicateUsername")
	ErrDuplicateEmail    = errors.New("identity.DuplicateEmail")
)

var (
	errEmptyUsername = errors.New("identity.validation.emptyUsername")
	errEmptyEmail    = errors.New("identity.validation.emptyEmail")
	errInvalidEmail  = errors.New("identity.validation.invalidEmail")
)
