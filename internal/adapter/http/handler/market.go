package handler

import (
	"encoding/json"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"marketflow/pkg/validator"
	"net/http"
	"time"
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
		log.Error("failed to get latest data from all exchanges", "error", err)
		internalErrorResponse(w, "failed to get latest data from all exchanges")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Error("failed to encode json", "err", err)
		internalErrorResponse(w, "failed to encode json")
		return
	}
}

// LatestPriceByExchange returns latest price for a specific exchange
func (h *Market) LatestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")

	log := h.log.GetSlogLogger().With("symbol", symbol)

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
		log.Error("failed to get latest data from specific exchange", "error", err)
		internalErrorResponse(w, "failed to get latest data from all exchanges")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Error("failed to encode json", "err", err)
		internalErrorResponse(w, "failed to encode json")
		return
	}
}

// HIGHEST

// HighestPrice returns highest price among all exchanges
func (h *Market) HighestPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")
	periodParsed, err := time.ParseDuration(period)
	if err != nil {
		h.log.Error(ctx, "failed to parse period, invalid format", "error", err)
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	result, err := h.market.GetHighest(ctx, "", symbol, periodParsed)

	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		internalErrorResponse(w, "failed to encode json")
		return
	}
}

// HighestPriceByExchange returns highest price for a sprcific exchange
func (h *Market) HighestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")
	periodParsed, err := time.ParseDuration(period)
	if err != nil {
		h.log.Error(ctx, "failed to parse period, invalid format", "error", err)
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	result, err := h.market.GetHighest(ctx, exchange, symbol, periodParsed)

	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		internalErrorResponse(w, "failed to encode json")
		return
	}
}

// LOWEST

// LowestPrice returns lowest price among all exchanges
func (h *Market) LowestPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")
	periodParsed, err := time.ParseDuration(period)
	if err != nil {
		h.log.Error(ctx, "failed to parse period, invalid format", "error", err)
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	result, err := h.market.GetLowest(ctx, "", symbol, periodParsed)

	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		internalErrorResponse(w, "failed to encode json")
		return
	}
}

// LowestPriceByExchange returns lowest price for a specific exchange
func (h *Market) LowestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")
	periodParsed, err := time.ParseDuration(period)
	if err != nil {
		h.log.Error(ctx, "failed to parse period, invalid format", "error", err)
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	result, err := h.market.GetLowest(ctx, exchange, symbol, periodParsed)

	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		internalErrorResponse(w, "failed to encode json")
		return
	}
}

// AVERAGE

// AveragePrice returns avg price among all exchanages
func (h *Market) AveragePrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")
	periodParsed, err := time.ParseDuration(period)
	if err != nil {
		h.log.Error(ctx, "failed to parse period, invalid format", "error", err)
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	result, err := h.market.GetAverage(ctx, "", symbol, periodParsed)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		internalErrorResponse(w, "failed to encode json")
		return
	}
}

// AveragePriceByExchange returns avg price for a specific exchange
func (h *Market) AveragePriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")
	period := r.URL.Query().Get("period")
	periodParsed, err := time.ParseDuration(period)
	if err != nil {
		h.log.Error(ctx, "failed to parse period, invalid format", "error", err)
		http.Error(w, "invalid period format", http.StatusBadRequest)
		return
	}

	result, err := h.market.GetAverage(ctx, exchange, symbol, periodParsed)

	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		internalErrorResponse(w, "failed to encode json")
		return
	}
}
