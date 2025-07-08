package types

import (
	"fmt"
	"slices"
)

// Symbol represents blockhain symbol/pairs
type Symbol string

var (
	BTCUSDT  Symbol = "BTCUSDT"
	DOGEUSDT Symbol = "DOGEUSDT"
	TONUSDT  Symbol = "TONUSDT"
	SOLUSDT  Symbol = "SOLUSDT"
	ETHUSDT  Symbol = "ETHUSDT"

	ValidSymbols = []Symbol{
		BTCUSDT, DOGEUSDT, TONUSDT, SOLUSDT, ETHUSDT,
	}

	ErrInvalidSymbol = fmt.Errorf("invalid symbol. Available: %v", ValidSymbols)
)

func (s Symbol) IsValid() bool {
	return slices.Contains(ValidSymbols, s)
}

func IsValidSymbol(s string) bool {
	return slices.Contains(ValidSymbols, Symbol(s))
}
