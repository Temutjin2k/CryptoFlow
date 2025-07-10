package postgres

import (
	"context"
	"marketflow/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	supportedExchanges = []string{"exchange1", "exchange2", "exchange3"}
)

type MarketRepo struct {
	db *pgxpool.Pool
}

func NewMarketRepository(db *pgxpool.Pool) *MarketRepo {
	return &MarketRepo{db: db}
}

// StoreStat inserts data stat to database
func (r *MarketRepo) StoreStats(ctx context.Context, stat *domain.PriceStats) error {
	query := `INSERT INTO aggregated_prices (pair_name, exchange, timestamp, min_price, max_price, average_price)
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(ctx, query,
		stat.Pair,
		stat.Exchange,
		stat.Timestamp,
		stat.Min,
		stat.Max,
		stat.Average)

	if err != nil {
		return ErrQueryFailed
	}

	return nil
}

// GetStats returns stats for given pair & exchange since given time
func (r *MarketRepo) GetStats(ctx context.Context, pair, exchange string, since time.Time) ([]*domain.PriceStats, error) {
	var query string
	var args []interface{}

	if exchange != "" {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1 AND exchange = $2 AND timestamp >= $3
		         ORDER BY timestamp DESC`
		args = []interface{}{pair, exchange, since}
	} else {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1 AND timestamp >= $2
		         ORDER BY timestamp DESC`
		args = []interface{}{pair, since}
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
	var args []interface{}

	if exchange != "" {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1 AND exchange = $2
		         ORDER BY timestamp DESC LIMIT 1`
		args = []interface{}{pair, exchange}
	} else {
		query = `SELECT pair_name, exchange, timestamp, min_price, max_price, average_price
		         FROM aggregated_prices
		         WHERE pair_name = $1
		         ORDER BY timestamp DESC LIMIT 1`
		args = []interface{}{pair}
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
