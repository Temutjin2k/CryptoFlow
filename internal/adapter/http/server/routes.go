package server

import (
	"net/http"
)

// setupRoutes - setups http routes
func (a *API) setupRoutes() {
	a.router.HandleFunc("/health", a.HealthCheck)

	Middleware(a.router)
}

func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) {

}
