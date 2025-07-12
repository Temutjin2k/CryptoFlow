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

var (
	StatusAvailable          = "available"
	StatusPartiallyAvailable = "partially_available"
	StatusHealthy            = "healthy"
	StatusUnhealthy          = "unhealthy"
)

// HealthCheck returns system information and service health statuses
func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Collect health status of all services
	serviceStatuses := make(map[string]any)
	healthyServices := 0
	totalServices := len(a.services)

	for _, svc := range a.services {
		healthy, err := svc.Health(r.Context())
		serviceStatus := StatusHealthy
		if err != nil {
			a.log.Error(r.Context(), "service health check failed",
				"service", svc.Name(), "error", err)
			serviceStatus = StatusUnhealthy
		} else if !healthy {
			serviceStatus = StatusUnhealthy
		} else {
			healthyServices++
		}

		serviceStatuses[svc.Name()] = serviceStatus
	}

	// Determine overall status
	var status string
	var statusCode int

	switch {
	case healthyServices == totalServices:
		status = StatusAvailable
		statusCode = http.StatusOK // 200
	case healthyServices > 0:
		status = StatusPartiallyAvailable
		statusCode = http.StatusPartialContent // 206
	default:
		status = StatusUnhealthy
		statusCode = http.StatusServiceUnavailable // 503
	}

	// Prepare response
	response := map[string]any{
		"system_info": map[string]any{
			"address":       a.addr,
			"data_mode":     a.modeProvider.Mode(),
			"status":        status,
			"total_healthy": healthyServices,
			"total":         totalServices,
			"services":      serviceStatuses,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.log.Error(r.Context(), "failed to encode health check response", "err", err)
		return
	}
}
