package service

import (
	"marketflow/internal/domain"
	"marketflow/internal/ports"
	"time"
)

type Market struct {
	repo ports.MarketRepository
}

func NewMarket(repo ports.MarketRepository) *Market {
	return &Market{repo: repo}
}

func (s *Market) GetLatest(exchange, symbol string) (domain.PriceData, error) {
	return domain.PriceData{}, nil
}

func (s *Market) GetHighest(exchange, symbol string, period time.Duration) (domain.PriceData, error) {
	return domain.PriceData{}, nil
}

func (s *Market) GetLowest(exchange, symbol string, period time.Duration) (domain.PriceData, error) {
	return domain.PriceData{}, nil
}

func (s *Market) GetAverage(exchange, symbol string, period time.Duration) (domain.PriceData, error) {
	return domain.PriceData{}, nil
}
