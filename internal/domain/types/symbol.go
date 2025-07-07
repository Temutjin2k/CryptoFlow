package types

import (
	"errors"
	"slices"
)

// Symbol represents blockhain symbol
type Symbol string

var (
	BTCUSDT  Symbol = "BTCUSDT"
	DOGEUSDT Symbol = "DOGEUSDT"
	TONUSDT  Symbol = "TONUSDT"
	SOLUSDT  Symbol = "SOLUSDT"
	ETHUSDT  Symbol = "ETHUSDT"

	// ValidSymbols = map[Symbol]struct{}{
	// 	BTCUSDT:  {},
	// 	DOGEUSDT: {},
	// 	TONUSDT:  {},
	// 	SOLUSDT:  {},
	// 	ETHUSDT:  {},
	// }

	ValidSymbols = []Symbol{
		BTCUSDT, DOGEUSDT, TONUSDT, SOLUSDT, ETHUSDT,
	}

	ErrInvalidSymbol = errors.New("invalid symbol")
)

func (s Symbol) IsValid() bool {
	return slices.Contains(ValidSymbols, s)
}

func IsValidSymbol(s string) bool {
	return slices.Contains(ValidSymbols, Symbol(s))
}
