package identity

import (
	"context"
	"net/mail"
)

func ValidateUser(ctx context.Context, user *User) error {
	if user.Username == "" {
		return errEmptyUsername
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
