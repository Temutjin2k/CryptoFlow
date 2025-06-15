package app

import (
	"marketflow/config"
	httpserver "marketflow/internal/adapter/http/server"
	"marketflow/pkg/logger"
	"marketflow/pkg/postgres"

	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "marketflow"

type App struct {
	httpServer *httpserver.API
	postgresDB *postgres.PostgreDB

	log logger.Logger
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

		log: logger,
	}
	return app, nil
}

func (app *App) Close(ctx context.Context) {

	// Closing database connection
	app.postgresDB.Pool.Close()

	// Closing http server
	err := app.httpServer.Stop()
	if err != nil {
		app.log.Info(ctx, "failed to shutdown HTTP service", "Err", err.Error())
	}

}

func (app *App) Run() error {
	errCh := make(chan error, 1)
	ctx := context.Background()

	// Running http server
	app.httpServer.Run(errCh)

	app.log.Info(ctx, "application started", "name", serviceName)

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun
	case s := <-shutdownCh:
		app.log.Info(ctx, "shuting down application", "signal", s.String())

		app.Close(ctx)
		app.log.Info(ctx, "graceful shutdown completed!")
	}

	return nil
}
