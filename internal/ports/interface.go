package ports

import (
	"context"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"time"
)

// redis
type Cache interface {
	SetLatest(ctx context.Context, latest domain.PriceData, duration time.Duration) error
	GetLatest(ctx context.Context, exchange types.Exchange, symbol types.Symbol) (*domain.PriceData, error)
}

// postgres
type MarketRepository interface {
	StoreStats(stat domain.PriceStats) error
	StoreStatsBatch(stats []domain.PriceStats) error
	GetStats(pair, exchange string, since time.Time) ([]*domain.PriceStats, error)
	GetLatest(ctx context.Context, exchange, pair string) (*domain.PriceData, error)
	GetByPeriod(ctx context.Context, exchange, pair string, period time.Duration) ([]*domain.PriceStats, error)
}

// ExchangeClient is an interface for Data sources
type ExchangeClient interface {
	Start(ctx context.Context) (<-chan domain.PriceData, error)
	Stop() error
}

type Collector interface {
	Start(processedPrices <-chan domain.PriceData)
}

// Service
type Market interface {
	// GetLatest returns latest price data from cache.
	GetLatest(ctx context.Context, exchange types.Exchange, symbol types.Symbol) (*domain.PriceData, error)
	GetHighest(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceData, error)
	GetLowest(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceData, error)
	GetAverage(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceData, error)
}
