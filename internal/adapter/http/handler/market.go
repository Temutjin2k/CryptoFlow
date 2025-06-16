package handler

import (
	"encoding/json"
	"marketflow/internal/adapter/http/service"
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

// GET localhost:8080/prices/latest/{symbol}
func (h *Market) LatestPrice(w http.ResponseWriter, r *http.Request) { //latest price over a period
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

// GET localhost:8080/prices/latest/{exchange}/{symbol}
func (h *Market) LatestPriceByExchange(w http.ResponseWriter, r *http.Request) { //latest price over a period by exchange
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

// GET localhost:8080/prices/highest/{symbol}
func (h *Market) HighestPrice(w http.ResponseWriter, r *http.Request) { //highest price over a period
	ctx := r.Context()
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	symbol := pathParts[3]
	period := r.URL.Query().Get("period")

	if period != "" {
		h.HighestPriceDuration(w, r)
		return
	}

	result, err := h.service.GetHighestPrice(symbol)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET localhost:8080/prices/highest/{exchange}/{symbol}
func (h *Market) HighestPriceByExchange(w http.ResponseWriter, r *http.Request) { //highest price over a period by exchange
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

	if period != "" {
		h.HighestPriceDurationByExchange(w, r)
		return
	}

	result, err := h.service.GetHighestPriceByExchange(exchange, symbol)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET localhost:8080/prices/highest/{symbol}?period={duration}
func (h *Market) HighestPriceDuration(w http.ResponseWriter, r *http.Request) { //highest price within a duration
	ctx := r.Context()
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	symbol := pathParts[3]
	period := r.URL.Query().Get("period")
	result, err := h.service.GetHighestPriceDuration(symbol, period)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET localhost:8080/prices/highest/{exchange}/{symbol}?period={duration}
func (h *Market) HighestPriceDurationByExchange(w http.ResponseWriter, r *http.Request) { //highest price within a duration by exchange
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
	result, err := h.service.GetHighestPriceDurationByExchange(exchange, symbol, period)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// LOWEST

// GET localhost:8080/prices/lowest/{symbol}
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

	if period != "" {
		h.LowestPriceDuration(w, r)
		return
	}

	result, err := h.service.GetLowestPrice(symbol)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET localhost:8080/prices/lowest/{exchange}/{symbol}
func (h *Market) LowestPriceByExchange(w http.ResponseWriter, r *http.Request) { //lowest price by exchange
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

	if period != "" {
		h.LowestPriceDurationByExchange(w, r)
		return
	}

	result, err := h.service.GetLowestPriceByExchange(exchange, symbol)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET localhost:8080/prices/lowest/{symbol}?period={duration}
func (h *Market) LowestPriceDuration(w http.ResponseWriter, r *http.Request) { //lowest price within a duration
	ctx := r.Context()
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	symbol := pathParts[3]
	period := r.URL.Query().Get("period")
	result, err := h.service.GetLowestPriceDuration(symbol, period)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET localhost:8080/prices/lowest/{exchange}/{symbol}?period={duration}
func (h *Market) LowestPriceDurationByExchange(w http.ResponseWriter, r *http.Request) { //lowest price within a duration by exchange
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
	result, err := h.service.GetLowestPriceDurationByExchange(exchange, symbol, period)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// AVERAGE

// GET localhost:8080/prices/average/{symbol}
func (h *Market) AveragePrice(w http.ResponseWriter, r *http.Request) { //average price over a period
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

// GET localhost:8080/prices/average/{exchange}/{symbol}?period={duration}
func (h *Market) AveragePriceDurationByExchange(w http.ResponseWriter, r *http.Request) { //average price within a duration by exchange
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
	result, err := h.service.GetAveragePriceDurationByExchange(exchange, symbol, period)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GET localhost:8080/prices/average/{exchange}/{symbol}
func (h *Market) AveragePriceByExchange(w http.ResponseWriter, r *http.Request) { //average price by exchange
	ctx := r.Context()
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		h.log.Error(ctx, "Invalid URL", "status", http.StatusBadRequest)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	exchange := pathParts[3]
	symbol := pathParts[4]
	result, err := h.service.GetAveragePriceByExchange(exchange, symbol)
	if err != nil {
		h.log.Error(ctx, "Failed to fetch data", "status", http.StatusInternalServerError)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
