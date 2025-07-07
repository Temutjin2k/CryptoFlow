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
		Redis    Redis
	}

	// Servers config
	Server struct {
		HTTPServer HTTPServer
	}

	// HTTP service
	HTTPServer struct {
		Port int `env:"HTTP_PORT" default:"8080"`
	}

	Redis struct {
		Addr     string `env:"REDIS_ADDR"`
		Password string `env:"REDIS_PASSWORD"`
	}

	Exchanges struct {
		Exchange1_Port string `env:"EXCHANGE1_PORT" default:"40101"`
		Exchange1_Name string `env:"EXCHANGE1_NAME" default:"exchange1"`

		Exchange2_Port string `env:"EXCHANGE2_PORT" default:"40102"`
		Exchange2_Name string `env:"EXCHANGE2_NAME" default:"exchange2"`

		Exchange3_Port string `env:"EXCHANGE3_PORT" default:"40103"`
		Exchange3_Name string `env:"EXCHANGE3_NAME" default:"exchange3"`
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
