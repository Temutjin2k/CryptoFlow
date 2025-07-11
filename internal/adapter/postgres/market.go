package postgres

import (
	"context"
	"marketflow/internal/domain"
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

	// Обрабатываем все ответы от batched-инсертов
	for range stats {
		if _, err := br.Exec(); err != nil {
			return ErrQueryFailed
		}
	}

	return nil
}

// GetStats returns stats for given pair & exchange since given time
func (r *MarketRepo) GetStats(ctx context.Context, pair, exchange string, since time.Time) ([]*domain.PriceStats, error) {
	var query string
	var args []any

	if exchange != "" {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1 AND exchange = $2 AND timestamp >= $3
		         ORDER BY timestamp DESC`
		args = []any{pair, exchange, since}
	} else {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1 AND timestamp >= $2
		         ORDER BY timestamp DESC`
		args = []any{pair, since}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, ErrQueryFailed
	}
	defer rows.Close()

	//scanning data to model
	var stats []*domain.PriceStats
	for rows.Next() {
		var stat domain.PriceStats
		if err := rows.Scan(
			&stat.Pair,
			&stat.Exchange,
			&stat.Timestamp,
			&stat.Min,
			&stat.Max,
			&stat.Average,
		); err != nil {
			return nil, ErrScanFailed
		}
		stats = append(stats, &stat)
	}

	if len(stats) == 0 {
		return nil, ErrNotFound
	}

	return stats, nil
}

// GetAverageStat returns latest stat row for given pair
func (r *MarketRepo) GetAverageStat(ctx context.Context, pair, exchange string) (*domain.PriceStats, error) {
	var query string
	var args []any

	if exchange != "" {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1 AND exchange = $2
		         ORDER BY timestamp DESC LIMIT 1`
		args = []any{pair, exchange}
	} else {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1
		         ORDER BY timestamp DESC LIMIT 1`
		args = []any{pair}
	}

	var stat domain.PriceStats
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&stat.Pair,
		&stat.Exchange,
		&stat.Timestamp,
		&stat.Min,
		&stat.Max,
		&stat.Average,
	)
	if err != nil {
		return nil, ErrScanFailed
	}

	return &stat, nil
}
