package config

import (
	"log/slog"
	"marketflow/pkg/envcfg"
	"marketflow/pkg/postgres"
)

type (
	// Config
	Config struct {
		Server   Server
		Postgres postgres.Config

		Logger *slog.Logger
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

	return config, err
}
