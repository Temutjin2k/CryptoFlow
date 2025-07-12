package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const serviceName = "postgresql"

type PostgreDB struct {
	Pool     *pgxpool.Pool
	DBConfig *pgxpool.Config
}

type Config struct {
	Dsn          string `env:"POSTGRES_DSN"`
	MaxOpenConns int32  `env:"POSTGRES_MAX_OPEN_CONN" default:"25"`
	MaxIdleTime  string `env:"POSTGRES_MAX_IDLE_TIME" default:"15m"`
}

func New(ctx context.Context, config Config) (*PostgreDB, error) {
	dbConfig, err := pgxpool.ParseConfig(config.Dsn)
	if err != nil {
		return nil, err
	}

	// Setting maxOpenConns
	dbConfig.MaxConns = config.MaxOpenConns

	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	// Setting MaxConnIdleTime
	dbConfig.MaxConnIdleTime = duration

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgreDB{
		Pool:     pool,
		DBConfig: dbConfig,
	}, nil
}

// Name returns the name of the database service
func (db *PostgreDB) Name() string {
	return serviceName
}

// Health checks the database connection health using Ping
func (db *PostgreDB) Health(ctx context.Context) (bool, error) {
	if db.Pool == nil {
		return false, fmt.Errorf("database pool is not initialized")
	}

	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Simply ping the database
	if err := db.Pool.Ping(healthCtx); err != nil {
		return false, fmt.Errorf("database ping failed: %w", err)
	}

	return true, nil
}
