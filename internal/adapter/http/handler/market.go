package handler

import (
	"encoding/json"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"marketflow/pkg/validator"
	"net/http"
	"time"
)

type Market struct {
	market ports.Market
	log    logger.Logger

	exchanges []string
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
	v := validator.New()
	if validateSymbol(v, symbol); !v.Valid() {
		errorResponse(w, r, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetLatest(ctx, "", symbol)

	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// LatestPriceByExchange returns latest price for a specific exchange
func (h *Market) LatestPriceByExchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	exchange := r.PathValue("exchange")
	symbol := r.PathValue("symbol")

	v := validator.New()

	h.validateExchange(v, exchange)

	if validateSymbol(v, symbol); !v.Valid() {
		errorResponse(w, r, http.StatusUnprocessableEntity, v.Errors)
		return
	}

	result, err := h.market.GetLatest(ctx, exchange, symbol)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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
	json.NewEncoder(w).Encode(result)
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
	json.NewEncoder(w).Encode(result)
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
	json.NewEncoder(w).Encode(result)
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
	json.NewEncoder(w).Encode(result)
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
	json.NewEncoder(w).Encode(result)
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
	json.NewEncoder(w).Encode(result)
}
