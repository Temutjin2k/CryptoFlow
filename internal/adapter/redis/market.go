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
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}

	data := new(domain.PriceData)
	if err := json.Unmarshal([]byte(val), data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PriceData: %w", err)
	}

	return data, nil
}

func (c *Cache) GetPriceInPeriod(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) ([]float64, error) {
	key := fmt.Sprintf("history:%s:%s", exchange, symbol)
	now := time.Now()
	start := now.Add(-period).UnixMilli()

	values, err := c.client.ZRangeByScore(ctx, key, &goredis.ZRangeBy{
		Min: fmt.Sprintf("%d", start),
		Max: fmt.Sprintf("%d", now.UnixMilli()),
	}).Result()

	if err != nil {
		return nil, err
	}
	prices := make([]float64, 0, len(values))
	for _, v := range values {
		var data domain.PriceData
		if err := json.Unmarshal([]byte(v), &data); err != nil {
			continue // skip bad record
		}
		prices = append(prices, data.Price)
	}
	return prices, nil
}

func (c *Cache) StoreHistory(ctx context.Context, p *domain.PriceData) error {
	key := fmt.Sprintf("history:%s:%s", p.Exchange, p.Symbol)
	score := float64(p.Timestamp.UnixMilli())
	value, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal PriceData: %w", err)
	}

	return c.client.ZAdd(ctx, key, goredis.Z{
		Score:  score,
		Member: value,
	}).Err()
}

// DeleteExpiredHistory deletes keys in history.* sorted sets that older than 5 minutes
func (c *Cache) DeleteExpiredHistory(ctx context.Context) error {
	retentionPeriod := c.cfg.HistoryDeleteDuration
	cutoff := time.Now().Add(-retentionPeriod).UnixMilli()

	// getting all keys by that has history:*
	keys, err := c.client.Keys(ctx, "history:*").Result()
	if err != nil {
		return fmt.Errorf("failed to get history keys: %w", err)
	}

	for _, key := range keys {
		err := c.client.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", cutoff)).Err()
		if err != nil {
			continue
		}
	}

	return nil
}
