package domain

import (
	"errors"
)

var (
	ErrUnimplemented    = errors.New("unimplemented")
	ErrInvalidSymbol    = errors.New("invalid symbol")
	ErrInvalidExchange  = errors.New("invalid exchange")
	ErrNegativePrice    = errors.New("price cannot be negative")
	ErrInvalidTimestamp = errors.New("invalid timestamp (zero time)")
)
