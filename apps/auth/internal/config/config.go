package config

import (
	"github.com/vandad1901/p3s/packages/go/envutil"
)

type Config struct {
	Environment envutil.Environment

	JWTConfig *JWTConfig
	DSN       string

	GRPCListenAddress string
	HTTPListenAddress string
}

const (
	gRPCListenAddress = "0.0.0.0:50051"
	HTTPListenAddress = "0.0.0.0:50151"
)

func LoadConfig() *Config {
	cfg := new(Config)
	cfg.Environment = envutil.MustGetEnvironment("APP_ENV")

	cfg.JWTConfig = loadJWTConfig()
	cfg.DSN = getDSN()

	switch cfg.Environment {
	case envutil.Development, envutil.Production:
		cfg.GRPCListenAddress = gRPCListenAddress
		cfg.HTTPListenAddress = HTTPListenAddress
	case envutil.Test:
	}

	return cfg
}
