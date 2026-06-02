package obviousutil

import "errors"

var obviousError = &obviousErr{}

type obviousErr struct {
	Key string
}

func (e *obviousErr) Error() string {
	return e.Key
}

func IsObvious(err error) bool {
	return errors.Is(err, obviousError)
}
