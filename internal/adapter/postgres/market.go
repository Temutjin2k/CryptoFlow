package postgres

import (
	"context"
	"errors"
	"fmt"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MarketRepo struct {
	db *pgxpool.Pool
}

func NewMarketRepository(db *pgxpool.Pool) *MarketRepo {
	return &MarketRepo{db: db}
}

// StoreStat inserts batch price stats to database
func (r *MarketRepo) StoreStats(ctx context.Context, stats []*domain.PriceStats) error {
	if len(stats) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	for _, stat := range stats {
		batch.Queue(`
			INSERT INTO aggregated_prices 
				(pair_name, exchange, timestamp, min_price, max_price, average_price)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			stat.Pair,
			stat.Exchange,
			stat.Timestamp,
			stat.Min,
			stat.Max,
			stat.Average,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for range stats {
		if _, err := br.Exec(); err != nil {
			return ErrQueryFailed
		}
	}

	return nil
}

// GetHighestStat returns highest price across exchanges in given period
func (r *MarketRepo) GetHighestStat(ctx context.Context, exchange types.Exchange, pair types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	timeThreshold := time.Now().Add(-period)

	var query string
	var args []any

	if exchange == types.AllExchanges {
		query = `
            SELECT 
                $1::text as pair_name,
                'ALL' as exchange,
                MAX(max_price) as max_price,
                MAX(timestamp) as timestamp
            FROM aggregated_prices
            WHERE pair_name = $1
            AND timestamp >= $2`
		args = []any{pair, timeThreshold}
	} else {
		query = `
            SELECT 
                $1::text as pair_name,
                $2::text as exchange,
                max_price,
                timestamp
            FROM aggregated_prices
            WHERE pair_name = $1
            AND exchange = $2
            AND timestamp >= $3
            ORDER BY max_price DESC
            LIMIT 1`
		args = []any{pair, exchange, timeThreshold}
	}

	var stats domain.PriceStats
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stats.Pair,
		&stats.Exchange,
		&stats.Max,
		&stats.Timestamp,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get highest stat: %w", err)
	}

	return &stats, nil
}

// GetLowestStat returns lowest price across exchanges in given period
func (r *MarketRepo) GetLowestStat(ctx context.Context, exchange types.Exchange, pair types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	timeThreshold := time.Now().Add(-period)

	var query string
	var args []any

	if exchange == types.AllExchanges {
		query = `
            SELECT 
                $1::text as pair_name,
                'ALL' as exchange,
                MIN(min_price) as min_price,
                MAX(timestamp) as timestamp
            FROM aggregated_prices
            WHERE pair_name = $1
            AND timestamp >= $2`
		args = []any{pair, timeThreshold}
	} else {
		query = `
            SELECT 
                $1::text as pair_name,
                $2::text as exchange,
                min_price,
                timestamp
            FROM aggregated_prices
            WHERE pair_name = $1
            AND exchange = $2
            AND timestamp >= $3
            ORDER BY min_price ASC
            LIMIT 1`
		args = []any{pair, exchange, timeThreshold}
	}

	var stats domain.PriceStats
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stats.Pair,
		&stats.Exchange,
		&stats.Min,
		&stats.Timestamp,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lowest stat: %w", err)
	}

	return &stats, nil
}

// GetAverageStat returns avarage price accross exchanges in given period
func (r *MarketRepo) GetAverageStat(ctx context.Context, exchange types.Exchange, pair types.Symbol, period time.Duration) (*domain.PriceStats, error) {
	timeThreshold := time.Now().Add(-period)

	var query string
	var args []any

	if exchange == types.AllExchanges {
		query = `
            SELECT 
                $1::text as pair_name,
                'ALL' as exchange,
                AVG(average_price) as average_price,
                MAX(timestamp) as timestamp
            FROM aggregated_prices
            WHERE pair_name = $1
            AND timestamp >= $2`
		args = []any{pair, timeThreshold}
	} else {
		query = `
            SELECT 
                $1::text as pair_name,
                $2::text as exchange,
                AVG(average_price) as average_price,
                MAX(timestamp) as timestamp
            FROM aggregated_prices
            WHERE pair_name = $1
            AND exchange = $2
            AND timestamp >= $3`
		args = []any{pair, exchange, timeThreshold}
	}

	var stats domain.PriceStats
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stats.Pair,
		&stats.Exchange,
		&stats.Average,
		&stats.Timestamp,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get average stat: %w", err)
	}

	return &stats, nil
}
