package postgres

import (
	"context"
	"marketflow/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MarketRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *MarketRepo {
	return &MarketRepo{db: db}
}

func (r *MarketRepo) StoreStats(stat domain.PriceStats) error {
	return domain.ErrUnimplemented
}

func (r *MarketRepo) StoreStatsBatch(stats []domain.PriceStats) error {
	return domain.ErrUnimplemented
}

func (r *MarketRepo) GetStats(pair, exchange string, since time.Time) ([]*domain.PriceStats, error) {
	return nil, domain.ErrUnimplemented
}

func (r *MarketRepo) GetLatest(ctx context.Context, exchange, pair string) (*domain.PriceStats, error) {
	return nil, domain.ErrUnimplemented
}

func (r *MarketRepo) GetByPeriod(ctx context.Context, exchange, pair string, period time.Duration) ([]*domain.PriceStats, error) {
	return nil, domain.ErrUnimplemented
}
