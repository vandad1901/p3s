package envutil

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const (
	EnvKeyDBHost = "PG_HOST"
	EnvKeyDBPort = "PG_PORT"
	EnvKeyDBUser = "PG_USER"
	EnvKeyDBPass = "PG_PASSWORD"
	EnvKeyDBName = "PG_DATABASE"

	EnvKeyJWTSecret = "JWT_SECRET"
)

var (
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	JWTSecret string
)

type ErrEnvVarNotSet struct {
	Key string
}

func (e *ErrEnvVarNotSet) Error() string {
	return "environment variable not set: " + e.Key
}

func load() {
	var found bool

	DBHost, found = os.LookupEnv(EnvKeyDBHost)
	if !found {
		panic(&ErrEnvVarNotSet{Key: EnvKeyDBHost})
	}

	DBPort, found = os.LookupEnv(EnvKeyDBPort)
	if !found {
		panic(&ErrEnvVarNotSet{Key: EnvKeyDBPort})
	}

	DBUser, found = os.LookupEnv(EnvKeyDBUser)
	if !found {
		panic(&ErrEnvVarNotSet{Key: EnvKeyDBUser})
	}

	DBPass, found = os.LookupEnv(EnvKeyDBPass)
	if !found {
		panic(&ErrEnvVarNotSet{Key: EnvKeyDBPass})
	}

	DBName, found = os.LookupEnv(EnvKeyDBName)
	if !found {
		panic(&ErrEnvVarNotSet{Key: EnvKeyDBName})
	}

	JWTSecret, found = os.LookupEnv(EnvKeyJWTSecret)
	if !found {
		panic(&ErrEnvVarNotSet{Key: EnvKeyJWTSecret})
	}
}
