package domain

import (
	"encoding/json"
	"fmt"
	"marketflow/internal/domain/types"
	"strings"
	"time"
)

type PriceData struct {
	Symbol    types.Symbol   `json:"symbol"`
	Exchange  types.Exchange `json:"exchange"`
	Price     float64        `json:"price"`
	Timestamp time.Time      `json:"timestamp"`
}

func (p *PriceData) IsValid() (bool, error) {
	if !p.Symbol.IsValid() {
		return false, ErrInvalidSymbol
	}
	if !p.Exchange.IsValid() {
		return false, ErrInvalidExchange
	}
	if p.Price < 0 {
		return false, ErrNegativePrice
	}
	if p.Timestamp.IsZero() {
		return false, ErrInvalidTimestamp
	}
	return true, nil
}

func (p PriceData) String() string {
	return fmt.Sprintf("[%s] %s = %.4f @ %s", p.Exchange, p.Symbol, p.Price, p.Timestamp.Format(time.RFC3339))
}

// Custom UnmarshalJSON to handle multiple timestamp formats
func (pd *PriceData) UnmarshalJSON(data []byte) error {
	type Alias PriceData
	aux := &struct {
		Timestamp interface{} `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(pd),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch v := aux.Timestamp.(type) {
	case string:
		// Handle quoted string timestamps
		t, err := time.Parse(time.RFC3339Nano, strings.Trim(v, `"`))
		if err != nil {
			return fmt.Errorf("invalid timestamp format: %w", err)
		}
		pd.Timestamp = t
	case float64:
		// Handle numeric timestamps (assuming milliseconds)
		pd.Timestamp = time.UnixMilli(int64(v))
	case json.Number:
		// Handle json.Number format
		ms, err := v.Int64()
		if err != nil {
			return fmt.Errorf("invalid numeric timestamp: %w", err)
		}
		pd.Timestamp = time.UnixMilli(ms)
	default:
		return fmt.Errorf("unsupported timestamp type: %T", v)
	}

	return nil
}

type PriceStats struct {
	Exchange  string
	Pair      string
	Timestamp time.Time
	Average   float64
	Min       float64
	Max       float64
}
