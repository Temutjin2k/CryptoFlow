package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type MarketRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *MarketRepo {
	return &MarketRepo{db: db}
}
