package handler

import (
	"fmt"
	"marketflow/internal/domain/types"
	"marketflow/pkg/validator"
)

func validateExchange(v *validator.Validator, exchange string) {
	v.Check(exchange != "", "exchange", "must be provided")

	v.Check(types.IsValidExchange(exchange), "exchange", ErrInvalidExchange)
}

func validateSymbol(v *validator.Validator, symbol string) {
	v.Check(symbol != "", "symbol", "must be provided")
	v.Check(types.IsValidSymbol(symbol), "symbol", ErrInvalidSymbol)
}

var (
	ErrInvalidExchange = fmt.Sprintf("invalid exchange. Available exchanges %v", types.ValidExchanges)
	ErrInvalidSymbol   = fmt.Sprintf("invalid symbol. Available: %v", types.ValidSymbols)
)
