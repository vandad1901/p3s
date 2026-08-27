package config

import (
	"log"

	"github.com/vandad1901/p3s/packages/go/envutil"
)

type Environment string

const (
	Development Environment = "development"
	Test        Environment = "test"
	Production  Environment = "production"
)

type Config struct {
	GRPCListenAddress string
	HTTPListenAddress string
	Environment       Environment

	JWTConfig *JWTConfig
	DSN       string
}

func LoadConfig() *Config {
	environment := envutil.MustGetString("APP_ENV")

	switch Environment(environment) {
	case Development, Test, Production:
	default:
		log.Fatalf("Invalid APP_ENV: %s", environment)
	}

	return &Config{
		GRPCListenAddress: envutil.MustGetString("GRPC_LISTEN_ADDRESS"),
		HTTPListenAddress: envutil.MustGetString("HTTP_LISTEN_ADDRESS"),
		Environment:       Environment(environment),

		JWTConfig: loadJWTConfig(),
		DSN:       getDSN(),
	}
}
