package app

import (
	"marketflow/config"
	"marketflow/internal/adapter/exchange"
	httpserver "marketflow/internal/adapter/http/server"
	repo "marketflow/internal/adapter/postgres"
	"marketflow/internal/adapter/redis"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
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

type App struct {
	httpServer      *httpserver.API
	postgresDB      *postgres.PostgreDB
	redis           *redis.Cache
	exchangeManager ports.ExchangeManager
	scheduler       ports.Sheduler

	log logger.Logger
}

func NewApplication(ctx context.Context, config config.Config, logger logger.Logger) (*App, error) {
	const fn = "app.NewApplication"

	log := logger.GetSlogLogger().With("fn", fn)
	// Postgres database
	db, err := postgres.New(ctx, config.Postgres)
	if err != nil {
		log.Error("failed to connect postgres", "dsn", config.Postgres.Dsn, "error", err)
		return nil, fmt.Errorf("failed to connect postgres: %v", err)
	}

	// Redis client
	cache, err := redis.NewClient(ctx, config.Redis)
	if err != nil {
		log.Error("failed to connect postgres", "address", config.Redis.Addr, "error", err)
		return nil, fmt.Errorf("failed to connect redis: %v", err)
	}

	// Define data sources
	exchange1 := exchange.NewExchange(types.Exchange1, config.DataManager.Exchanges.Exchange1Addr, logger)
	exchange2 := exchange.NewExchange(types.Exchange2, config.DataManager.Exchanges.Exchange2Addr, logger)
	exchange3 := exchange.NewExchange(types.Exchange3, config.DataManager.Exchanges.Exchange3Addr, logger)

	// Test mode data sources
	// exchange1 := exchange.GenerateTestData()
	// exchange2 := exchange.GenerateTestData()
	// exchange3 := exchange.GenerateTestData()

	sources := []ports.ExchangeSource{
		exchange1,
		exchange2,
		exchange3,
	}

	//repository
	repo := repo.NewMarketRepository(db.Pool)

	// Market service
	market := service.NewMarket(repo, cache, logger)

	// Aggregator
	aggregator := service.NewAggregator(repo, cache, config.DataManager.Aggregator.TickerDuration, logger)

	// DataCollector
	collector := service.NewCollector(cache, nil, logger)

	// ExchangeManager
	exchangeManager := service.NewExchangeManager(sources, collector, aggregator, config.DataManager.Distributor.WorkerCount, logger)

	scheduler := service.NewScheduler(ctx, logger)
	scheduler.AddTask("Delete expired exchange history", types.TaskTypeInterval, config.Redis.HistoryDeleteDuration, cache.DeleteExpiredHistory)

	// List of all services for healthcheck
	serviceList := []httpserver.Service{
		exchange1,
		exchange2,
		exchange3,
		db,
		cache,
	}

	// REST API server
	httpServer := httpserver.New(config, market, serviceList, logger)

	app := &App{
		httpServer:      httpServer,
		postgresDB:      db,
		exchangeManager: exchangeManager,
		redis:           cache,
		scheduler:       scheduler,
		log:             logger,
	}
	return app, nil
}

func (app *App) close(ctx context.Context) {
	app.scheduler.Close()

	// Closing http server
	if err := app.httpServer.Stop(); err != nil {
		app.log.Warn(ctx, "failed to shutdown HTTP service", "error", err)
	}

	if err := app.exchangeManager.Close(); err != nil {
		app.log.Warn(ctx, "failed to shutdown exchange manager", "error", err)
	}

	// Closing database connection
	app.postgresDB.Pool.Close()

	// Closing redis
	app.redis.Close()

}

func (app *App) Run() error {
	const fn = "app.Run"
	log := app.log.GetSlogLogger().With("fn", fn)

	errCh := make(chan error, 1)
	ctx := context.Background()
	// Running DataManager
	if err := app.exchangeManager.Start(ctx); err != nil {
		log.Error("failed to start exchange manager", "error", err)
		return err
	}

	// Running http server
	app.httpServer.Run(errCh)

	// Running scheduler
	app.scheduler.Start()

	log.InfoContext(ctx, "application started", "name", serviceName)

	// Waiting signal
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case errRun := <-errCh:
		return errRun
	case s := <-shutdownCh:
		log.InfoContext(ctx, "shuting down application", "signal", s.String())

		app.close(ctx)
		log.InfoContext(ctx, "graceful shutdown completed!")
	}

	return nil
}
