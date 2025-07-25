package config

import (
	"fmt"
	"time"

	"marketflow/pkg/envcfg"
	"marketflow/pkg/loadenv"
	"marketflow/pkg/postgres"
)

type (
	// Config
	Config struct {
		Server      Server
		Postgres    postgres.Config
		Redis       Redis
		DataManager DataManager
	}

	Test struct {
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
		Addr                  string        `env:"REDIS_ADDR"`
		Password              string        `env:"REDIS_PASSWORD"`
		HistoryDeleteDuration time.Duration `env:"REDIS_HISTORY_DELETE_DURATION" default:"5m"`
	}

	DataManager struct {
		Exchanges   Exchanges
		Distributor Distributor
		Aggregator  Aggregator
	}

	Distributor struct {
		WorkerCount int `env:"DISTRIBUTOR_WORKER_COUNT" default:"5"`
	}

	// Exchanges config
	Exchanges struct {
		Exchange1Addr string `env:"EXCHANGE1_ADDR" default:"localhost:40101"`
		Exchange2Addr string `env:"EXCHANGE2_ADDR" default:"localhost:40102"`
		Exchange3Addr string `env:"EXCHANGE3_ADDR" default:"localhost:40103"`
	}

	Aggregator struct {
		TickerDuration time.Duration `env:"AGGREGATOR_TICKER_DURATION" default:"1m"`
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
