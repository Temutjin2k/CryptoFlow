package domain

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrUnimplemented    = errors.New("unimplemented")
	ErrInvalidSymbol    = errors.New("invalid symbol")
	ErrInvalidExchange  = errors.New("invalid exchange")
	ErrNegativePrice    = errors.New("price cannot be negative")
	ErrInvalidTimestamp = errors.New("invalid timestamp (zero time)")

	ErrAlreadyOnLiveMode = errors.New("server is already on live mode")
	ErrAlreadyOnTestMode = errors.New("server is already on test mode")
)
