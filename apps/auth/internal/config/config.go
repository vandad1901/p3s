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
}

func LoadConfig() *Config {
	gRPCListenAddress := envutil.MustGetString("GRPC_LISTEN_ADDRESS")
	httpListenAddress := envutil.MustGetString("HTTP_LISTEN_ADDRESS")
	environment := envutil.MustGetString("APP_ENV")

	switch Environment(environment) {
	case Development, Test, Production:
	default:
		log.Fatalf("Invalid APP_ENV: %s", environment)
	}

	return &Config{
		GRPCListenAddress: gRPCListenAddress,
		HTTPListenAddress: httpListenAddress,
		Environment:       Environment(environment),
	}
}
