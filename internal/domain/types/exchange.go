package types

import (
	"slices"
)

// Exchange defines type to represent Exchanges(data sources).
type Exchange string

var (
	AllExchanges Exchange = "all"
	Exchange1    Exchange = "exchange1"
	Exchange2    Exchange = "exchange2"
	Exchange3    Exchange = "exchange3"

	ValidExchanges = []Exchange{
		Exchange1,
		Exchange2,
		Exchange3,
	}
)

func (e *Exchange) IsValid() bool {
	return IsValidExchange(string(*e))
}

func IsValidExchange(s string) bool {
	return slices.Contains(ValidExchanges, Exchange(s))
}
