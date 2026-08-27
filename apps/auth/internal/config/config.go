package config

import (
	"github.com/vandad1901/p3s/packages/go/envutil"
)

type Config struct {
	GRPCListenAddress string
	HTTPListenAddress string
	Environment       envutil.Environment

	JWTConfig *JWTConfig
	DSN       string
}

func LoadConfig() *Config {
	environment := envutil.MustGetEnvironment("APP_ENV")

	var (
		gRPCListenAddress string
		HTTPListenAddress string
	)

	switch environment {
	case envutil.Development, envutil.Production:
		gRPCListenAddress = envutil.MustGetString("GRPC_LISTEN_ADDRESS")
		HTTPListenAddress = envutil.MustGetString("HTTP_LISTEN_ADDRESS")
	case envutil.Test:
	}

	return &Config{
		GRPCListenAddress: gRPCListenAddress,
		HTTPListenAddress: HTTPListenAddress,
		Environment:       environment,

		JWTConfig: loadJWTConfig(),
		DSN:       getDSN(),
	}
}
