package domain

import (
	"fmt"
	"marketflow/internal/domain/types"
	"time"
)

type PriceData struct {
	Symbol    types.Symbol `json:"symbol"`
	Exchange  string       `json:"exchange"`
	Price     float64      `json:"price"`
	Timestamp time.Time    `json:"timestamp"`
}

func (p PriceData) String() string {
	return fmt.Sprintf("[%s] %s = %.4f @ %s", p.Exchange, p.Symbol, p.Price, p.Timestamp.Format(time.RFC3339))
}

type PriceStats struct {
	Exchange  string
	Pair      string
	Timestamp time.Time
	Average   float64
	Min       float64
	Max       float64
}
