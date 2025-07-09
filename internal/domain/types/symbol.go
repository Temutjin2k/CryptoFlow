package types

import (
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
)

func (s Symbol) IsValid() bool {
	return slices.Contains(ValidSymbols, s)
}

func IsValidSymbol(s string) bool {
	return slices.Contains(ValidSymbols, Symbol(s))
}
