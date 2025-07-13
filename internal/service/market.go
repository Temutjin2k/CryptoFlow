package service

import (
	"context"
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
)

type Market struct {
	storage ports.MarketRepository
	cache   ports.Cache
	logger  logger.Logger
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

	if latest == nil {
		return nil, domain.ErrNotFound
	}

	return latest, nil
}

func (s *Market) GetHighest(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	const fn = "GetHighest"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	if period < time.Minute {
		return s.fetchHighestFromCache(ctx, exchange, symbol, period)
	}

	highest, err := s.storage.GetHighestStat(ctx, exchange, symbol, period)
	if err != nil {
		log.Error("failed to get stats from database", "error", err)
		return s.fetchHighestFromCache(ctx, exchange, symbol, period)
	}

	if highest == nil {
		return s.fetchHighestFromCache(ctx, exchange, symbol, period)
	}

	return &domain.PriceStats{
		Exchange:  highest.Exchange,
		Pair:      highest.Pair,
		Timestamp: highest.Timestamp,
		Max:       highest.Max,
	}, nil
}

func (s *Market) fetchHighestFromCache(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	prices, err := s.cache.GetPriceInPeriod(ctx, exchange, symbol, period)
	if err != nil {
		s.logger.Error(ctx, "failed to get data from Cache", "exchange", exchange, "symbol", symbol, "period", period, "error", err)
		return nil, nil
	}

	_, max, _ := aggregateAndPrice(prices)

	if len(prices) == 0 {
		s.logger.Warn(ctx, "no prices found in cache")
		return nil, domain.ErrNotFound
	}

	return &domain.PriceStats{
		Exchange:  exchange,
		Pair:      symbol,
		Timestamp: max.Timestamp,
		Max:       max.Price,
	}, nil
}

func (s *Market) GetLowest(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	const fn = "GetLowest"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	if period < time.Minute {
		return s.fetchLowestFromCache(ctx, exchange, symbol, period)
	}

	lowest, err := s.storage.GetLowestStat(ctx, exchange, symbol, period)
	if err != nil {
		log.Error("failed to get stats from database", "error", err)
		return s.fetchLowestFromCache(ctx, exchange, symbol, period)
	}

	if lowest == nil {
		return s.fetchLowestFromCache(ctx, exchange, symbol, period)
	}

	return &domain.PriceStats{
		Exchange:  lowest.Exchange,
		Pair:      lowest.Pair,
		Timestamp: lowest.Timestamp,
		Min:       lowest.Min,
	}, nil
}

func (s *Market) fetchLowestFromCache(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	prices, err := s.cache.GetPriceInPeriod(ctx, exchange, symbol, period)
	if err != nil {
		s.logger.Error(ctx, "failed to get data from Cache", "exchange", exchange, "symbol", symbol, "period", period, "error", err)
		return nil, nil
	}

	if len(prices) == 0 {
		s.logger.Warn(ctx, "no prices found in cache")
		return nil, domain.ErrNotFound
	}

	min, _, _ := aggregateAndPrice(prices)

	return &domain.PriceStats{
		Exchange:  exchange,
		Pair:      symbol,
		Timestamp: min.Timestamp,
		Min:       min.Price,
	}, nil
}

func (s *Market) GetAverage(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	const fn = "GetAverage"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	if period < time.Minute {
		return s.fetchAverageFromCache(ctx, exchange, symbol, period)
	}

	avg, err := s.storage.GetAverageStat(ctx, exchange, symbol, period)
	if err != nil {
		log.Error("failed to get stats from database, trying to check from cache...", "error", err)
		return s.fetchAverageFromCache(ctx, exchange, symbol, period)
	}

	if avg == nil {
		return s.fetchAverageFromCache(ctx, exchange, symbol, period)
	}

	return &domain.PriceStats{
		Exchange:  avg.Exchange,
		Pair:      avg.Pair,
		Timestamp: avg.Timestamp,
		Average:   avg.Average,
	}, nil
}

func (s *Market) fetchAverageFromCache(ctx context.Context, exchange types.Exchange, symbol types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	prices, err := s.cache.GetPriceInPeriod(ctx, exchange, symbol, period)
	if err != nil {
		s.logger.Error(ctx, "failed to get data from Cache", "exchange", exchange, "symbol", symbol, "period", period, "error", err)
		return nil, nil
	}

	if len(prices) == 0 {
		s.logger.Warn(ctx, "no prices found in cache")
		return nil, domain.ErrNotFound
	}

	_, _, avg := aggregateAndPrice(prices)

	return &domain.PriceStats{
		Exchange:  exchange,
		Pair:      symbol,
		Timestamp: avg.Timestamp,
		Average:   avg.Price,
	}, nil
}
