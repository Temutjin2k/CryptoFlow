package ports

import (
	"context"
	"marketflow/internal/domain"
	"time"
)

// redis
type Cache interface {
	SetLatest(ctx context.Context, update domain.PriceData) error
	GetLatest(ctx context.Context, exchange, pair string) (domain.PriceData, error)
}

// postgres
type MarketRepository interface {
	StoreStats(stat domain.PriceStats) error
	StoreStatsBatch(stats []domain.PriceStats) error
	GetStats(pair, exchange string, since time.Time) ([]domain.PriceStats, error)
	GetLatest(ctx context.Context, exchange, pair string) (domain.PriceStats, error)
	GetByPeriod(ctx context.Context, exchange, pair string, period time.Duration) ([]domain.PriceStats, error)
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
	GetLatest(exchange, symbol string) (domain.PriceData, error)
	GetHighest(exchange, symbol string, period time.Duration) (domain.PriceData, error)
	GetLowest(exchange, symbol string, period time.Duration) (domain.PriceData, error)
	GetAverage(exchange, symbol string, period time.Duration) (domain.PriceData, error)
}
