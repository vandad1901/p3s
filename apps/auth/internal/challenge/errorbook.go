package challenge

import "errors"

var (
	errEmptyType        = errors.New("challenge.validation.emptyType")
	errInvalidMetadata  = errors.New("challenge.validation.invalidMetadata")
	errExpiredChallenge = errors.New("challenge.consume.expired")
	errAlreadyConsumed  = errors.New("challenge.consume.alreadyConsumed")
	errBadHash          = errors.New("challenge.consume.badHash")
)
