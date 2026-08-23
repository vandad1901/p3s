package identity

import "errors"

var (
	errEmptyUsername = errors.New("identity.validation.emptyUsername")
	errEmptyEmail    = errors.New("identity.validation.emptyEmail")
	errInvalidEmail  = errors.New("identity.validation.invalidEmail")
)
