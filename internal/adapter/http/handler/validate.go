package handler

import (
	"marketflow/internal/domain/types"
	"marketflow/pkg/validator"
)

func validateExchange(v *validator.Validator, exchange string) {
	v.Check(exchange != "", "exchange", "must be provided")

	v.Check(types.IsValidExchange(exchange), "exchange", types.ErrInvalidExchange.Error())
}

func validateSymbol(v *validator.Validator, symbol string) {
	v.Check(symbol != "", "symbol", "must be provided")
	v.Check(types.IsValidSymbol(symbol), "symbol", types.ErrInvalidSymbol.Error())
}
