package identity

import (
	"context"
	"errors"
	"net/mail"
)

func ValidateUser(ctx context.Context, user *User) error {
	if user.Username == "" {
		return errors.New("identity.MandatoryUsername")
	}

	if user.Email == "" {
		return errors.New("identity.MandatoryEmail")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return errors.New("identity.InvalidEmail")
	}

	return nil
}
