package service

import (
	"context"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"time"
)

type Market struct {
	storage ports.MarketRepository
	cache   ports.Cache

	logger logger.Logger
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

func (s *Market) GetHighest(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceData, error) {
	return nil, domain.ErrUnimplemented
}

func (s *Market) GetLowest(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceData, error) {
	return nil, domain.ErrUnimplemented
}

func (s *Market) GetAverage(ctx context.Context, exchange, symbol string, period time.Duration) (*domain.PriceData, error) {
	return nil, domain.ErrUnimplemented
}
