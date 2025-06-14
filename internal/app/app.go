package app

import (
	"marketflow/config"
	httpserver "marketflow/internal/adapter/http/server"
	"marketflow/pkg/logger"
	"marketflow/pkg/postgres"

	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "marketflow"

type App struct {
	httpServer *httpserver.API
	postgresDB *postgres.PostgreDB

	log *slog.Logger
}

func NewApplication(ctx context.Context, config config.Config, logger logger.Logger) (*App, error) {
	// Database
	db, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres: %v", err)
	}

	httpServer := httpserver.New(config, logger)

	app := &App{
		httpServer: httpServer,
		postgresDB: db,

		log: config.Logger,
	}
	return app, nil
}

func (app *App) Close(ctx context.Context) {

	// Closing database connection
	app.postgresDB.Pool.Close()

	// Closing http server
	err := app.httpServer.Stop()
	if err != nil {
		app.log.Info("failed to shutdown HTTP service", "Err", err.Error())
	}

}

func (app *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	// Running http server
	app.httpServer.Run(errCh)

	app.log.Info("service started", "name", serviceName)

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun
	case s := <-shutdownCh:
		app.log.Info("Shuting down gracefully", "signal", s)

		app.Close(ctx)
		app.log.Info("graceful shutdown completed!")
	}

	return nil
}
