package config

import (
	"fmt"
	"marketflow/pkg/envcfg"
	"marketflow/pkg/loadenv"
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

	// Custom func that loads enviromental variables
	if err := loadenv.LoadEnvFile(".env"); err != nil {
		return Config{}, fmt.Errorf("failed to load enviromental variables: %w", err)
	}

	// Parsing enviromental variables to the struct
	if err := envcfg.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}
