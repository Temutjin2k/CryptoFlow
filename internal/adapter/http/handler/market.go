package handler

import (
	"errors"
	"net/http"

	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"marketflow/pkg/validator"
)

type Market struct {
	market ports.Market
	log    logger.Logger
}

func NewMarket(market ports.Market, log logger.Logger) *Market {
	return &Market{
		market: market,
		log:    log,
	}
}

// LATEST

// LatestPrice returns latest price among all exchanges
func (h *Market) LatestPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	symbol := r.PathValue("symbol")

	log := h.log.GetSlogLogger().With("symbol", symbol)

	v := validator.New()
	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	// getting latest price data from all exchanges
	result, err := h.market.GetLatest(ctx, types.AllExchanges, types.Symbol(symbol))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}
		log.Error("failed to get latest data from all exchanges", "error", err)
		internalErrorResponse(w, "failed to get latest data from all exchanges")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result}, nil)
}

// LatestPriceByExchange returns latest price for a specific exchange
func (h *Market) LatestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")

	log := h.log.GetSlogLogger().With("symbol", symbol, "exchange", exchange)

	v := validator.New()

	validateExchange(v, exchange)

	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	// getting latest price data from specific exchange.
	result, err := h.market.GetLatest(ctx, types.Exchange(exchange), types.Symbol(symbol))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}
		log.Error("failed to get latest data from specific exchange", "error", err)
		internalErrorResponse(w, "failed to get latest data from all exchanges")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result}, nil)
}

// HIGHEST

// HighestPrice returns highest price among all exchanges
func (h *Market) HighestPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")

	log := h.log.GetSlogLogger().With("symbol", symbol)

	v := validator.New()

	periodParsed, normalizedPeriod, err := parsePeriod(period)
	if err != nil {
		v.AddError("period", err.Error())
	}

	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetHighest(ctx, types.AllExchanges, types.Symbol(symbol), periodParsed)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}
		h.log.Error(ctx, "failed to fetch highest price", "error", err)
		internalErrorResponse(w, "failed to fetch highest price")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result, "period": normalizedPeriod}, nil)
}

// HighestPriceByExchange returns highest price for a sprcific exchange
func (h *Market) HighestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")

	log := h.log.GetSlogLogger().With("symbol", symbol, "exchange", exchange)

	v := validator.New()

	periodParsed, normalizedPeriod, err := parsePeriod(period)
	if err != nil {
		v.AddError("period", err.Error())
	}

	validateExchange(v, exchange)

	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetHighest(ctx, types.Exchange(exchange), types.Symbol(symbol), periodParsed)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}

		h.log.Error(ctx, "failed to fetch highest price by exchange", "error", err)
		internalErrorResponse(w, "failed to fetch highest price by exchange")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result, "period": normalizedPeriod}, nil)
}

// LOWEST

// LowestPrice returns lowest price among all exchanges
func (h *Market) LowestPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")

	log := h.log.GetSlogLogger().With("symbol", symbol)

	v := validator.New()

	periodParsed, normalizedPeriod, err := parsePeriod(period)
	if err != nil {
		v.AddError("period", err.Error())
	}

	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetLowest(ctx, types.AllExchanges, types.Symbol(symbol), periodParsed)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}

		h.log.Error(ctx, "failed to fetch lowest price", "error", err)
		internalErrorResponse(w, "failed to fetch lowest price")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result, "period": normalizedPeriod}, nil)
}

// LowestPriceByExchange returns lowest price for a specific exchange
func (h *Market) LowestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")

	log := h.log.GetSlogLogger().With("symbol", symbol, "exchange", exchange)

	v := validator.New()
	periodParsed, normalizedPeriod, err := parsePeriod(period)
	if err != nil {
		v.AddError("period", err.Error())
	}

	validateExchange(v, exchange)

	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetLowest(ctx, types.Exchange(exchange), types.Symbol(symbol), periodParsed)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}

		h.log.Error(ctx, "failed to fetch lowest price by exchange", "error", err)
		internalErrorResponse(w, "failed to fetch lowest price by exchange")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result, "period": normalizedPeriod}, nil)
}

// AVERAGE

// AveragePrice returns avg price among all exchanages
func (h *Market) AveragePrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")

	log := h.log.GetSlogLogger().With("symbol", symbol)

	v := validator.New()

	periodParsed, normalizedPeriod, err := parsePeriod(period)
	if err != nil {
		v.AddError("period", err.Error())
	}

	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetAverage(ctx, types.AllExchanges, types.Symbol(symbol), periodParsed)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}

		h.log.Error(ctx, "failed to fetch lowest average", "error", err)
		internalErrorResponse(w, "failed to fetch average price by exchange")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result, "period": normalizedPeriod}, nil)
}

// AveragePriceByExchange returns avg price for a specific exchange
func (h *Market) AveragePriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")

	log := h.log.GetSlogLogger().With("symbol", symbol, "exchange", exchange)

	v := validator.New()

	periodParsed, normalizedPeriod, err := parsePeriod(period)
	if err != nil {
		v.AddError("period", err.Error())
	}

	validateExchange(v, exchange)

	if validateSymbol(v, symbol); !v.Valid() {
		log.Error("failed to validate request", "errors", v.Errors)
		errorResponse(w, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetAverage(ctx, types.Exchange(exchange), types.Symbol(symbol), periodParsed)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			notFoundErrorResponse(w)
			return
		}

		h.log.Error(ctx, "failed to fetch average price by exchange", "error", err)
		internalErrorResponse(w, "failed to fetch average price by exchange")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"data": result, "period": normalizedPeriod}, nil)
}
