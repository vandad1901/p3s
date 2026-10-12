package identity

import (
	"context"
	"net/mail"
)

const (
	MinUsernameLength    = 5
	MinDisplayNameLength = 3
)

func ValidateUser(ctx context.Context, user *User) error {
	if len(user.Username) < MinUsernameLength {
		return errInvalidUsername
	}

	if len(user.DisplayName) < MinDisplayNameLength {
		return errInvalidDisplayName
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
