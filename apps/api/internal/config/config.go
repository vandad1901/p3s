package config

import (
	"github.com/vandad1901/p3s/packages/go/envutil"
)

type Config struct {
	Environment envutil.Environment

	DSN string

	GRPCListenAddress string
}

const (
	gRPCListenAddress = "0.0.0.0:50052"
)

func LoadConfig() *Config {
	cfg := new(Config)
	cfg.Environment = envutil.MustGetEnvironment("APP_ENV")

	cfg.DSN = getDSN()

	switch cfg.Environment {
	case envutil.Development, envutil.Production:
		cfg.GRPCListenAddress = gRPCListenAddress
	case envutil.Test:
	}

	return cfg
}
