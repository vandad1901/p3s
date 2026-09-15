package config

import (
	"github.com/vandad1901/p3s/packages/go/envutil"
)

type Config struct {
	Environment envutil.Environment

	RabbitMQAddress string

	S3Endpoint     string
	S3RootUsername string
	S3RootPassword string

	DSN string

	AuthServiceAddress string
}

func LoadConfig() *Config {
	cfg := new(Config)

	cfg.Environment = envutil.MustGetEnvironment("APP_ENV")

	cfg.RabbitMQAddress = envutil.MustGetString("RMQ_ENDPOINT")

	cfg.S3Endpoint = envutil.MustGetString("S3_ENDPOINT")
	cfg.S3RootUsername = envutil.MustGetString("S3_ROOT_USERNAME")
	cfg.S3RootPassword = envutil.MustGetString("S3_ROOT_PASSWORD")

	cfg.DSN = getDSN()

	cfg.AuthServiceAddress = envutil.MustGetString("AUTH_JWKS_ADDRESS")

	return cfg
}
