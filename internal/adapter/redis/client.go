package redis

import (
	"context"
	"fmt"
	"marketflow/config"
	"marketflow/internal/domain/types"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const serviceName = "redis"

type Cache struct {
	client *goredis.Client

	cfg config.Redis
}

// NewClient returns new Redis-client
func NewClient(ctx context.Context, cfg config.Redis) (*Cache, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})

	// Ping Redis to verify the connection
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := rdb.Ping(pingCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis at %s: %w", cfg.Addr, err)
	}
	return &Cache{
		client: rdb,
		cfg:    cfg,
	}, nil
}

// Closes the client
func (c *Cache) Close() error {
	if c.client == nil {
		return nil // Already closed or not initialized
	}

	// Gracefully close connections
	err := c.client.Close()
	c.client = nil // Prevent double closing

	if err != nil {
		return fmt.Errorf("redis close error: %w", err)
	}
	return nil
}

// Name returns the name of the Redis service
func (c *Cache) Name() string {
	return serviceName
}

// Health checks the Redis connection health using Ping
func (c *Cache) Health(ctx context.Context) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}

	// Create a context with timeout for the health check
	healthCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Use Ping to check Redis health
	if err := c.client.Ping(healthCtx).Err(); err != nil {
		return false, fmt.Errorf("redis ping failed: %w", err)
	}

	return true, nil
}

func (c *Cache) createKeyByExchangeAndSymbol(exchange types.Exchange, symbol types.Symbol) string {
	return fmt.Sprintf("latest:%s:%s", exchange, symbol)
}

func (c *Cache) createKeyBySymbol(symbol types.Symbol) string {
	return fmt.Sprintf("latest:%s", symbol)
}

func (c *Cache) createHistoryKeyByExchangeAndSymbol(exchange types.Exchange, symbol types.Symbol) string {
	return fmt.Sprintf("history:%s:%s", exchange, symbol)
}

func (c *Cache) createHistoryKeyBySymbol(symbol types.Symbol) string {
	return fmt.Sprintf("history:%s", symbol)
}
