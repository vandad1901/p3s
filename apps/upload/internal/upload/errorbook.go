package upload

import "errors"

var (
	errEmptyKey   = errors.New("upload.validation.emptyKey")
	errInvalidKey = errors.New("upload.validation.invalidKey")
)
