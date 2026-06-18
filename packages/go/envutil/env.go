package envutil

import (
	"os"
	"strconv"
)

const (
	notFoundError = "environment variable not found: "
)

func GetString(key string) (string, bool) {
	value, found := os.LookupEnv(key)

	return value, found
}

func MustGetString(key string) string {
	value, found := GetString(key)
	if !found {
		panic(notFoundError + key)
	}

	return value
}

func GetInt64(key string) (int64, bool) {
	value, found := os.LookupEnv(key)
	if !found {
		return 0, false
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}

	return intValue, true
}

func MustGetInt64(key string) int64 {
	value, found := GetInt64(key)
	if !found {
		panic(notFoundError + key)
	}

	return value
}

func GetBool(key string) (bool, bool) {
	value, found := os.LookupEnv(key)
	if !found {
		return false, false
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, false
	}

	return boolValue, true
}

func MustGetBool(key string) bool {
	value, found := GetBool(key)
	if !found {
		panic(notFoundError + key)
	}

	return value
}
