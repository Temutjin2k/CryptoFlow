package server

import (
	"encoding/json"
	"net/http"
)

// setupRoutes - setups http routes
func (a *API) setupRoutes() {
	// System Health
	a.router.HandleFunc("/health", a.HealthCheck)

	// Market Data API

	//latest
	a.router.HandleFunc("/prices/latest/{symbol}", a.routes.market.LatestPrice)
	a.router.HandleFunc("/prices/latest/{exchange}/{symbol}", a.routes.market.LatestPriceByExchange)

	//highest
	a.router.HandleFunc("/prices/highest/{symbol}", a.routes.market.HighestPrice)
	a.router.HandleFunc("/prices/highest/{exchange}/{symbol}", a.routes.market.HighestPriceByExchange)
	a.router.HandleFunc("/prices/highest/{symbol}?period={duration}", a.routes.market.HighestPriceDuration)
	a.router.HandleFunc("/prices/highest/{exchange}/{symbol}?period={duration}", a.routes.market.HighestPriceDurationByExchange)

	//lowest
	a.router.HandleFunc("/prices/lowest/{symbol}", a.routes.market.LowestPrice)
	a.router.HandleFunc("/prices/lowest/{exchange}/{symbol}", a.routes.market.LowestPriceByExchange)
	a.router.HandleFunc("/prices/lowest/{symbol}?period={duration}", a.routes.market.LowestPriceDuration)
	a.router.HandleFunc("/prices/lowest/{exchange}/{symbol}?period={duration}", a.routes.market.LowestPriceDurationByExchange)

	//average
	a.router.HandleFunc("/prices/average/{symbol}", a.routes.market.AveragePrice)
	a.router.HandleFunc("/prices/average/{exchange}/{symbol}", a.routes.market.AveragePriceByExchange)
	a.router.HandleFunc("/prices/average/{exchange}/{symbol}?period={duration}", a.routes.market.AveragePriceDurationByExchange)

	// Data Mode API
	a.router.HandleFunc("POST /mode/test", a.routes.mode.TestMode) // Switch to Test Mode (use generated data).
	a.router.HandleFunc("POST /mode/live", a.routes.mode.LiveMode) // Switch to Live Mode (fetch data from provided programs).
}

// HealthCheck - returns system information.
func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"address": a.addr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.log.Error(r.Context(), "failed to encode", "err", err)
		return
	}
}
