package config

import (
	"github.com/vandad1901/p3s/packages/go/envutil"
)

type Config struct {
	Environment envutil.Environment

	RabbitMQAddress string

	S3Endpoint         string
	S3ExternalEndpoint string
	S3RootUsername     string
	S3RootPassword     string

	DSN string

	AuthServiceAddress string

	GRPCListenAddress string
}

const GRPCListenAddress = "0.0.0.0:50053"

func LoadConfig() *Config {
	cfg := new(Config)
	cfg.Environment = envutil.MustGetEnvironment("APP_ENV")

	cfg.RabbitMQAddress = envutil.MustGetString("RMQ_ENDPOINT")

	cfg.S3Endpoint = envutil.MustGetString("S3_ENDPOINT")
	cfg.S3ExternalEndpoint = envutil.MustGetString("S3_EXTERNAL_ENDPOINT")
	cfg.S3RootUsername = envutil.MustGetString("S3_ROOT_USERNAME")
	cfg.S3RootPassword = envutil.MustGetString("S3_ROOT_PASSWORD")

	cfg.DSN = getDSN()

	switch cfg.Environment {
	case envutil.Development, envutil.Production:
		cfg.GRPCListenAddress = GRPCListenAddress
		cfg.AuthServiceAddress = envutil.MustGetString("AUTH_JWKS_ADDRESS")
	case envutil.Test:
	}

	return cfg
}
