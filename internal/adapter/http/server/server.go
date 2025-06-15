package server

import (
	"context"
	"errors"
	"fmt"
	"marketflow/config"
	"marketflow/internal/adapter/http/handler"
	"marketflow/pkg/logger"
	"net/http"
	"time"
)

const serverIPAddress = ":%d" // Change to 0.0.0.0 for external access

type API struct {
	cfg    config.HTTPServer
	router *http.ServeMux
	server *http.Server

	addr string

	routes *handlers // routes/handlers

	log logger.Logger
}

type handlers struct {
	someHandler handler.SomeHandler
}

func New(cfg config.Config, logger logger.Logger) *API {
	addr := fmt.Sprintf(serverIPAddress, cfg.Server.HTTPServer.Port)

	someHandler := handler.NewSomeHandler()

	handlers := &handlers{
		someHandler: *someHandler,
	}

	// Setup routes
	mux := http.NewServeMux()

	api := &API{
		router: mux,
		routes: handlers,

		addr: addr,
		cfg:  cfg.Server.HTTPServer,
		log:  logger,
	}

	api.server = &http.Server{
		Addr:    api.addr,
		Handler: api.router,
	}

	api.setupRoutes()

	return api
}

func (a *API) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.log.Info(ctx, "Shutting down HTTP server...", "Address", a.addr)
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	return nil
}

func (a *API) Run(errCh chan<- error) {
	go func() {
		a.log.Info(context.Background(), "Started http server", "Address", a.addr)
		if err := http.ListenAndServe(a.addr, a.applyMiddlewares(a.router)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to start HTTP server: %w", err)
			return
		}
	}()
}

// applyMiddlewares to wrap default http.ServeMux
func (m *API) applyMiddlewares(next http.Handler) http.Handler {
	return m.LoggingMiddleware(next)
}
