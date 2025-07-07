package service

import (
	"context"
	"marketflow/internal/domain"
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

func (s *Market) GetLatest(ctx context.Context, exchange, symbol string) (*domain.PriceData, error) {
	const fn = "GetLatest"
	log := s.logger.GetSlogLogger().With("fn", fn, "exchange", exchange, "symbol", symbol)

	latest, err := s.cache.GetLatest(ctx, exchange, symbol)
	if err != nil {
		log.Warn("failed to get latest data from cache, trying to check from storage...", "error", err)

		latest, err = s.storage.GetLatest(ctx, exchange, symbol)
		if err != nil {
			log.Error("failed to get latest data from storage", "error", err)
			return nil, err
		}
	}

	return latest, domain.ErrUnimplemented
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
