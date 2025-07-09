package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = errors.New("value with given key not found")
)

// SetLatest saves PriceData into Redis 2 keys(by exchange and symbol, and by symbol only) with given TTL(Time-To-Live)
func (c *Cache) SetLatest(ctx context.Context, latest *domain.PriceData, ttl time.Duration) error {
	key := c.createKeyByExchangeAndSymbol(latest.Exchange, latest.Symbol) // Key by exchange and symbol
	keyBySymbol := c.createKeyBySymbol(latest.Symbol)

	data, err := json.Marshal(latest)
	if err != nil {
		return fmt.Errorf("failed to marshal PriceData: %w", err)
	}

	// Setting key by exchange and symbol
	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	// Setting key only by symbol
	if err := c.client.Set(ctx, keyBySymbol, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return nil
}

// GetLatest returns PriceData by given key (exchange, pair)
func (c *Cache) GetLatest(ctx context.Context, exchange types.Exchange, symbol types.Symbol) (*domain.PriceData, error) {
	var key string
	if exchange == types.AllExchanges {
		key = c.createKeyBySymbol(symbol)
	} else {
		key = c.createKeyByExchangeAndSymbol(exchange, symbol)
	}

	val, err := c.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	data := new(domain.PriceData)
	if err := json.Unmarshal([]byte(val), data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PriceData: %w", err)
	}

	return data, nil
}
