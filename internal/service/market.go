package service

import (
	"context"
	"marketflow/internal/adapter/postgres"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"sync"
	"time"
)

type Market struct {
	storage ports.MarketRepository
	cache   ports.Cache
	logger  logger.Logger
	mu      sync.RWMutex
}

func NewMarket(repo ports.MarketRepository, cache ports.Cache, logger logger.Logger) *Market {
	return &Market{
		storage: repo,
		cache:   cache,
		logger:  logger,
	}
}

// GetLatest returns latest price data from cache.
func (s *Market) GetLatest(ctx context.Context, exchange types.Exchange, symbol types.Symbol) (*domain.PriceData, error) {
	const fn = "GetLatest"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	latest, err := s.cache.GetLatest(ctx, exchange, symbol)
	if err != nil {
		log.Error("failed to get latest data from cache", "error", err)
		return nil, err
	}

	return latest, nil
}

// AggregateAndStore gets aggregated data and sends it to database
func (s *Market) AggregateAndStore(ctx context.Context) {
	s.logger.Info(ctx, "Running AggregateAndStore")
	exchanges := types.ValidExchanges
	symbols := types.ValidSymbols

	for _, exchange := range exchanges {
		for _, symbol := range symbols {
			values, err := s.cache.GetPriceInPeriod(ctx, exchange, symbol, time.Minute)
			if len(values) == 0 {
				s.logger.Warn(ctx, "No prices found wait for 1 minute", "exchange", exchange, "symbol", symbol, "len", len(values), "err", err)
				continue
			}
			if err != nil {
				s.logger.Error(ctx, "failed to get prices from cache", "exchange", exchange, "symbol", symbol, "error", err)
				continue
			}

			min, max, avg := aggregate(values)

			stat := &domain.PriceStats{
				Exchange:  string(exchange),
				Pair:      string(symbol),
				Timestamp: time.Now().Truncate(time.Minute),
				Average:   avg,
				Min:       min,
				Max:       max,
			}

			//to prevent race condition
			// s.mu.Lock()
			err = s.storage.StoreStats(ctx, stat)
			// s.mu.Unlock()

			if err != nil {
				s.logger.Error(ctx, "failed to save stat", "exchange", exchange, "symbol", symbol, "error", err)
			}
		}
	}
}

func (s *Market) GetHighest(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceStats, error) {
	const fn = "GetHighest"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	if symbol == "" {
		return nil, domain.ErrInvalidSymbol
	}

	since := time.Now().Add(-period).Truncate(time.Minute)

	//to prevent race condition
	// s.mu.RLock()
	stats, err := s.storage.GetStats(ctx, symbol, exchange, since)
	// s.mu.RUnlock()

	if err != nil {
		log.Error("failed to get stats from database", "error", err)
		return nil, err
	}

	if len(stats) == 0 {
		return nil, postgres.ErrNotFound
	}

	highest := stats[0]
	for _, v := range stats[1:] {
		if v != nil && v.Max > highest.Max {
			highest = v
		}
	}

	return &domain.PriceStats{
		Exchange:  highest.Exchange,
		Pair:      highest.Pair,
		Timestamp: highest.Timestamp,
		Max:       highest.Max,
	}, nil
}

func (s *Market) GetLowest(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceStats, error) {
	const fn = "GetLowest"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	since := time.Now().Add(-period - time.Minute)

	s.mu.RLock()
	stats, err := s.storage.GetStats(ctx, symbol, exchange, since)
	s.mu.RUnlock()

	if len(stats) == 0 {
		return nil, postgres.ErrNotFound
	}

	if err != nil {
		log.Error("failed to get stats from database", "error", err)
		return nil, err
	}

	lowest := stats[0]
	for _, v := range stats[1:] {
		if v.Max < lowest.Min {
			lowest = v
		}
	}

	return &domain.PriceStats{
		Exchange:  lowest.Exchange,
		Pair:      lowest.Pair,
		Timestamp: lowest.Timestamp,
		Min:       lowest.Min,
	}, nil
}

func (s *Market) GetAverage(ctx context.Context, exchange, symbol string) (*domain.PriceStats, error) {
	const fn = "GetAverage"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	avg, err := s.storage.GetAverageStat(ctx, symbol, exchange)
	if err != nil {
		log.Error("failed to get stats from database", "error", err)
		return nil, err
	}

	if avg == nil {
		return nil, postgres.ErrNotFound
	}

	return &domain.PriceStats{
		Exchange:  avg.Exchange,
		Pair:      avg.Pair,
		Timestamp: avg.Timestamp,
		Average:   avg.Average,
	}, nil
}

func aggregate(values []float64) (min, max, avg float64) {
	min, max = values[0], values[0]
	sum := 0.0

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}

	avg = sum / float64(len(values))
	return
}
