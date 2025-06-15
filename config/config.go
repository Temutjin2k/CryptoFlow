package config

import (
	"marketflow/pkg/envcfg"
	"marketflow/pkg/postgres"
)

type (
	// Config
	Config struct {
		Server   Server
		Postgres postgres.Config
	}

	// Servers config
	Server struct {
		HTTPServer HTTPServer
	}

	// HTTP service
	HTTPServer struct {
		Port int `env:"HTTP_PORT" default:"8080"`
	}
)

func New() (Config, error) {
	var config Config

	err := envcfg.Parse(&config)

	if config.Postgres.Dsn == "" {
		config.Postgres.Dsn = "postgres://user:admin@localhost:5432/marketflow?sslmode=disable"
	}

	return config, err
}
