package envutil

import "log"

type Environment string

const (
	Development Environment = "development"
	Test        Environment = "test"
	Production  Environment = "production"
)

func GetEnvironment(key string) (Environment, bool) {
	value, found := GetString(key)
	if !found {
		return "", false
	}

	switch Environment(value) {
	case Test, Development, Production:
	default:
		log.Fatalf("Invalid APP_ENV: %s", value)
	}

	return Environment(value), true
}

func MustGetEnvironment(key string) Environment {
	value, found := GetEnvironment(key)
	if !found {
		panic(notFoundError + key)
	}

	return value
}
