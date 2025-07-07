package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketflow/config"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("value with given key not found")
)

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

// SetLatest saves PriceData into Redis with given TTL(Time-To-Live)
func (c *Cache) SetLatest(ctx context.Context, latest domain.PriceData, ttl time.Duration) error {
	key := c.createKeyByExchangeAndSymbol(latest.Exchange, latest.Symbol)

	data, err := json.Marshal(latest)
	if err != nil {
		return fmt.Errorf("failed to marshal PriceData: %w", err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// GetLatest returns PriceData by given key (exchange, pair)
func (c *Cache) GetLatest(ctx context.Context, exchange, pair string) (*domain.PriceData, error) {
	key := fmt.Sprintf("latest:%s:%s", exchange, pair)

	val, err := c.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	var data *domain.PriceData
	if err := json.Unmarshal([]byte(val), data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PriceData: %w", err)
	}

	return data, nil
}

func (c *Cache) createKeyByExchangeAndSymbol(exchange string, symbol types.Symbol) string {
	return fmt.Sprintf("latest:%s:%s", exchange, symbol)
}
