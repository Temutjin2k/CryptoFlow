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
	// Latest
	a.router.HandleFunc("/prices/latest/{symbol}", a.routes.market.LatestPrice)
	a.router.HandleFunc("/prices/latest/{exchange}/{symbol}", a.routes.market.LatestPriceByExchange)

	// Highest
	a.router.HandleFunc("/prices/highest/{symbol}", a.routes.market.HighestPrice)
	a.router.HandleFunc("/prices/highest/{exchange}/{symbol}", a.routes.market.HighestPriceByExchange)

	// Lowest
	a.router.HandleFunc("/prices/lowest/{symbol}", a.routes.market.LowestPrice)
	a.router.HandleFunc("/prices/lowest/{exchange}/{symbol}", a.routes.market.LowestPriceByExchange)

	// Average
	a.router.HandleFunc("/prices/average/{symbol}", a.routes.market.AveragePrice)
	a.router.HandleFunc("/prices/average/{exchange}/{symbol}", a.routes.market.AveragePriceByExchange)

	// Data Mode
	a.router.HandleFunc("/mode/test", a.routes.mode.TestMode)
	a.router.HandleFunc("/mode/live", a.routes.mode.LiveMode)
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
