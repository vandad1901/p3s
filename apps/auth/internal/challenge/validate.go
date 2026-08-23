package challenge

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

func validateOnCreate(_ context.Context, _ *gorm.DB, ch *Challenge) error {
	err := validateType(ch.ChallengeType)
	if err != nil {
		return err
	}

	err = validateMetadata(ch.Metadata)
	if err != nil {
		return err
	}

	return nil
}

func validateType(challengeType string) error {
	if challengeType == "" {
		return errEmptyType
	}

	return nil
}

func validateMetadata(metadata string) error {
	ok := json.Valid([]byte(metadata))
	if !ok {
		return errInvalidMetadata
	}

	return nil
}
