package handler

import (
	"encoding/json"
	"marketflow/internal/adapter/http/service"
	model "marketflow/internal/domain"
	"marketflow/pkg/logger"
	"net/http"
	"strings"
)

type Market struct {
	service service.MarketService
	log     logger.Logger
}

func NewMarket(log logger.Logger) *Market {
	return &Market{
		log: log,
	}
}

// LATEST

// LatestPrice returns latest price among all exchanges
func (h *Market) LatestPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) < 4 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	symbol := pathParts[3]
	result, err := h.service.GetLatestPrice(symbol)

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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	exchange := pathParts[3]
	symbol := pathParts[4]
	result, err := h.service.GetLatestPriceByExchange(exchange, symbol)
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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	symbol := pathParts[3]
	period := r.URL.Query().Get("period")
	var result model.PriceData
	var err error

	if period == "" {
		result, err = h.service.GetHighestPrice(symbol)
	} else {
		result, err = h.service.GetHighestPriceDuration(symbol, period)
	}

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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	exchange := pathParts[3]
	symbol := pathParts[4]
	period := r.URL.Query().Get("period")

	var result model.PriceData
	var err error

	if period == "" {
		result, err = h.service.GetHighestPriceByExchange(exchange, symbol)
	} else {
		result, err = h.service.GetHighestPriceDurationByExchange(exchange, symbol, period)
	}

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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	symbol := pathParts[3]
	period := r.URL.Query().Get("period")

	var result model.PriceData
	var err error

	if period == "" {
		result, err = h.service.GetLowestPrice(symbol)
	} else {
		result, err = h.service.GetLowestPriceDuration(symbol, period)
	}

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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	exchange := pathParts[3]
	symbol := pathParts[4]
	period := r.URL.Query().Get("period")

	var result model.PriceData
	var err error

	if period == "" {
		result, err = h.service.GetLowestPriceByExchange(exchange, symbol)
	} else {
		result, err = h.service.GetLowestPriceDurationByExchange(exchange, symbol, period)
	}

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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	symbol := pathParts[3]
	result, err := h.service.GetAveragePrice(symbol)
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
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	exchange := pathParts[3]
	symbol := pathParts[4]
	period := r.URL.Query().Get("period")
	var result model.PriceData
	var err error

	if period == "" {
		result, err = h.service.GetAveragePriceByExchange(exchange, symbol)
	} else {
		result, err = h.service.GetAveragePriceDurationByExchange(exchange, symbol, period)
	}

	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
