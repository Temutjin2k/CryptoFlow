package server

import (
	"encoding/json"
	"net/http"
)

// setupRoutes - setups http routes
func (a *API) setupRoutes() {
	a.router.HandleFunc("/health", a.HealthCheck)
}

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
