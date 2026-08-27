package config

import (
	"github.com/vandad1901/p3s/packages/go/envutil"
)

type Config struct {
	GRPCListenAddress string
	Environment       envutil.Environment

	DSN string
}

func LoadConfig() *Config {
	environment := envutil.MustGetEnvironment("APP_ENV")

	var gRPCListenAddress string

	switch environment {
	case envutil.Development, envutil.Production:
		gRPCListenAddress = envutil.MustGetString("GRPC_LISTEN_ADDRESS")
	case envutil.Test:
	}

	return &Config{
		GRPCListenAddress: gRPCListenAddress,
		Environment:       environment,

		DSN: getDSN(),
	}
}
