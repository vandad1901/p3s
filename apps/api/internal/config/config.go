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
	Environment       Environment

	DSN string
}

func LoadConfig() *Config {
	environment := envutil.MustGetString("APP_ENV")

	var gRPCListenAddress string

	switch Environment(environment) {
	case Development, Production:
		gRPCListenAddress = envutil.MustGetString("GRPC_LISTEN_ADDRESS")
	case Test:
	default:
		log.Fatalf("Invalid APP_ENV: %s", environment)
	}

	return &Config{
		GRPCListenAddress: gRPCListenAddress,
		Environment:       Environment(environment),

		DSN: getDSN(),
	}
}
