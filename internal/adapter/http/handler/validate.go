package handler

import (
	"fmt"
	"marketflow/internal/domain/types"
	"marketflow/pkg/validator"
)

func (h *Market) validateExchange(v *validator.Validator, exchange string) {
	v.Check(exchange != "", "exchange", "must be provided")

	if h.exchanges != nil {
		v.Check(validator.PermittedValue(exchange, h.exchanges...), "exchange", fmt.Sprintf("invalid exchange. Available exchanges %v", h.exchanges))
	}
}

func validateSymbol(v *validator.Validator, symbol string) {
	v.Check(symbol != "", "symbol", "must be provided")
	v.Check(types.IsValidSymbol(symbol), "symbol", fmt.Sprintf("invalid symbol. Available symbols %v", types.ValidSymbols))
}
