package ports

import (
	"context"
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
)

// redis
type Cache interface {
	SetLatest(ctx context.Context, latest *domain.PriceData, duration time.Duration) error
	GetLatest(ctx context.Context, exchange types.Exchange, symbol types.Symbol) (*domain.PriceData, error)
	GetPriceInPeriod(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) ([]*domain.PriceData, error)
	StoreHistory(ctx context.Context, p *domain.PriceData) error
}

type ExchangeManager interface {
	Start(ctx context.Context) error
	Close() error
}

// postgres
type MarketRepository interface {
	StoreStats(ctx context.Context, stat []*domain.PriceStats) error
	GetHighestStat(ctx context.Context, exchange types.Exchange, pair types.Symbol, period time.Duration) (*domain.PriceStats, error)
	GetAverageStat(ctx context.Context, exchange types.Exchange, pair types.Symbol, period time.Duration) (*domain.PriceStats, error)
	GetLowestStat(ctx context.Context, exchange types.Exchange, pair types.Symbol, period time.Duration) (*domain.PriceStats, error)
}

// ExchangeSource is an interface for Data sources
type ExchangeSource interface {
	Name() string
	Start(ctx context.Context) (<-chan *domain.PriceData, error)
	Close() error
}

type Distributor interface {
	FanOut(ctx context.Context)
}

type WorkerPool interface {
	Start(ctx context.Context)
	Input() chan<- *domain.PriceData
	Output() <-chan *domain.PriceData
	Close()
}

type Aggregator interface {
	Start(ctx context.Context)
	FanIn(ctx context.Context, inputs ...<-chan *domain.PriceData) <-chan *domain.PriceData
}

type Collector interface {
	Start(ctx context.Context, processedPrices <-chan *domain.PriceData)
	Cancel() error
}

type Sheduler interface {
	Start()
	Close()
	AddTask(name string, taskType types.TaskType, interval time.Duration, handler func(ctx context.Context) error)
}

// Service
type Market interface {
	// GetLatest returns latest price data from cache.
	GetLatest(ctx context.Context, exchange types.Exchange, symbol types.Symbol) (*domain.PriceData, error)
	GetHighest(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error)
	GetLowest(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error)
	GetAverage(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error)
}
