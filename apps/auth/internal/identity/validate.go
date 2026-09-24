package identity

import (
	"context"
	"net/mail"
)

const (
	MinUsernameLength = 5
)

func ValidateUser(ctx context.Context, user *User) error {
	if user.Username == "" {
		return errEmptyUsername
	}

	if len(user.Username) < MinUsernameLength {
		return errInvalidUsername
	}

	if user.Email == "" {
		return errEmptyEmail
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errInvalidEmail
	}

	return nil
}
