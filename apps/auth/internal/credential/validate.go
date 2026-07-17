package credential

import "errors"

const PASSWORD_MINIMUM_LENGTH = 12

func ValidatePassword(password string) error {
	// We haven't used utf8.RuneCount on purpose, as the pure byte count is more important in regards to entropy
	if len(password) < PASSWORD_MINIMUM_LENGTH {
		return errors.New("credential.InsecurePassword")
	}

	// TODO: add validation via haveibeenpwned.com

	return nil
}
