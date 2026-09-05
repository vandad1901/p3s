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

	HTTPListenAddress string
}

func LoadConfig() *Config {
	environment := envutil.MustGetEnvironment("APP_ENV")

	var HTTPListenAddress string

	switch environment {
	case envutil.Development, envutil.Production:
		HTTPListenAddress = envutil.MustGetString("HTTP_LISTEN_ADDRESS")
	case envutil.Test:
	}

	return &Config{
		Environment: environment,

		RabbitMQAddress: envutil.MustGetString("UPLOAD_RMQ_ENDPOINT"),

		S3Endpoint:     envutil.MustGetString("UPLOAD_S3_ENDPOINT"),
		S3RootUsername: envutil.MustGetString("UPLOAD_S3_ROOT_USERNAME"),
		S3RootPassword: envutil.MustGetString("UPLOAD_S3_ROOT_PASSWORD"),
		S3BucketName:   envutil.MustGetString("UPLOAD_S3_BUCKET_NAME"),

		DSN: getDSN(),

		AuthServiceAddress: envutil.MustGetString("AUTH_JWKS_ADDRESS"),

		HTTPListenAddress: HTTPListenAddress,
	}
}
