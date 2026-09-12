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
	S3BucketName   string

	DSN string

	AuthServiceAddress string
}

func LoadConfig() *Config {
	environment := envutil.MustGetEnvironment("APP_ENV")

	return &Config{
		Environment: environment,

		RabbitMQAddress: envutil.MustGetString("MEDIA_RMQ_ENDPOINT"),

		S3Endpoint:     envutil.MustGetString("MEDIA_S3_ENDPOINT"),
		S3RootUsername: envutil.MustGetString("MEDIA_S3_ROOT_USERNAME"),
		S3RootPassword: envutil.MustGetString("MEDIA_S3_ROOT_PASSWORD"),
		S3BucketName:   envutil.MustGetString("MEDIA_S3_BUCKET_NAME"),

		DSN: getDSN(),

		AuthServiceAddress: envutil.MustGetString("AUTH_JWKS_ADDRESS"),
	}
}
