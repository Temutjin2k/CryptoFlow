package app

import (
	"marketflow/config"
	httpserver "marketflow/internal/adapter/http/server"
	"marketflow/internal/adapter/redis"
	"marketflow/internal/service"
	"marketflow/pkg/logger"
	"marketflow/pkg/postgres"

	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "marketflow"

type ExchangeManager interface {
	Start(ctx context.Context) error
}

type App struct {
	httpServer      *httpserver.API
	postgresDB      *postgres.PostgreDB
	exchangeManager ExchangeManager

	log logger.Logger
}

func NewApplication(ctx context.Context, config config.Config, logger logger.Logger) (*App, error) {
	const fn = "NewApplication"

	log := logger.GetSlogLogger().With("fn", fn)
	// Postgres database
	db, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		log.Error("failed to connect postgres", "dsn", config.Postgres.Dsn, "error", err)
		return nil, fmt.Errorf("failed to connect postgres: %v", err)
	}

	// Redis client
	redisClient, err := redis.NewClient(ctx, config.Redis)
	if err != nil {
		log.Error("failed to connect postgres", "adress", config.Redis.Addr, "error", err)
		return nil, fmt.Errorf("failed to connect redis: %v", err)
	}
	// // Define data sources
	// sources := []ports.ExchangeClient{
	// 	exchange.NewExchange("exchange1", "localhost:40101", logger),
	// 	exchange.NewExchange("exchange2", "localhost:40101", logger),
	// 	exchange.NewExchange("exchange3", "localhost:40101", logger),
	// }

	// DataCollector
	// collector := service.NewCollector(nil, nil, logger)

	// // ExchangeManager
	// exchangeManager := service.NewExchangeManager(collector, sources, logger)

	// Market service
	market := service.NewMarket(nil, redisClient, logger)

	// REST API server
	httpServer := httpserver.New(config, market, logger)

	app := &App{
		httpServer:      httpServer,
		postgresDB:      db,
		exchangeManager: nil,

		log: logger,
	}
	return app, nil
}

func (app *App) close(ctx context.Context) {
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

	// Running DataManager
	// if err := app.exchangeManager.Start(ctx); err != nil {
	// 	return err
	// }

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

		app.close(ctx)
		app.log.Info(ctx, "graceful shutdown completed!")
	}

	return nil
}
