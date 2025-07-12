package service

import (
	"context"
	"fmt"

	"marketflow/config"
	"marketflow/internal/adapter/exchange"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
)

// ExchangeManager manages all working process related to exchanges
type ExchangeManager struct {
	exchangeSources []ports.ExchangeSource
	distributors    []ports.Distributor
	workerPools     []ports.WorkerPool
	aggregator      ports.Aggregator
	collector       ports.Collector

	store ports.MarketRepository
	cache ports.Cache

	isLive bool

	cfg    config.DataManager
	logger logger.Logger
}

func NewExchangeManager(
	islive bool,
	exchanges []ports.ExchangeSource,
	store ports.MarketRepository,
	cache ports.Cache,

	cfg config.DataManager,
	logger logger.Logger,
) *ExchangeManager {
	return &ExchangeManager{
		exchangeSources: exchanges,
		store:           store,
		cache:           cache,

		isLive: islive,
		cfg:    cfg,
		logger: logger,
	}
}

func (m *ExchangeManager) Start(ctx context.Context) error {
	m.initCollectorAndAggregator()

	for _, source := range m.exchangeSources {
		m.logger.Info(ctx, "starting source", "name", source.Name())
		pricesCh, err := source.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start source: %w", err)
		}

		workerPool := NewWorkerPool(source.Name(), m.cfg.Distributor.WorkerCount, m.logger)
		m.workerPools = append(m.workerPools, workerPool)

		distributor := NewDistriubtor(workerPool, pricesCh)
		m.distributors = append(m.distributors, distributor)

		workerPool.Start(ctx)
		distributor.FanOut(ctx)
	}

	merged := m.aggregator.FanIn(ctx, m.getWorkerPoolOutputs()...)

	// starting collector
	m.collector.Start(ctx, merged)

	// starting aggregator
	m.aggregator.Start(ctx)

	return nil
}

func (m *ExchangeManager) Close() error {
	const fn = "ExchangeManager.Close"
	log := m.logger.GetSlogLogger().With("fn", fn)

	for _, source := range m.exchangeSources {
		if err := source.Close(); err != nil {
			log.Warn("failed to close exchange", "error", err)
		}
	}

	if err := m.collector.Cancel(); err != nil {
		log.Warn("failed to cancel collector", "error", err)
	}

	for _, pool := range m.workerPools {
		pool.Close()
	}

	return nil
}

func (m *ExchangeManager) SwitchToTest() error {
	if !m.isLive {
		return domain.ErrAlreadyOnTestMode
	}
	m.Close()

	testSources := []ports.ExchangeSource{
		exchange.NewTestExchange(types.Exchange1),
		exchange.NewTestExchange(types.Exchange2),
		exchange.NewTestExchange(types.Exchange3),
	}
	// switching to test sources
	m.exchangeSources = testSources

	m.isLive = false

	return m.Start(context.Background())
}

func (m *ExchangeManager) SwitchToLive() error {
	if m.isLive {
		return domain.ErrAlreadyOnLiveMode
	}
	m.Close()

	liveSources := []ports.ExchangeSource{
		exchange.NewExchange(types.Exchange1, m.cfg.Exchanges.Exchange1Addr, m.logger),
		exchange.NewExchange(types.Exchange2, m.cfg.Exchanges.Exchange2Addr, m.logger),
		exchange.NewExchange(types.Exchange3, m.cfg.Exchanges.Exchange3Addr, m.logger),
	}
	// switching to live sources
	m.exchangeSources = liveSources

	m.isLive = true

	return m.Start(context.Background())
}

func (m *ExchangeManager) Mode() string {
	if m.isLive {
		return types.LiveMode
	}

	return types.TestMode
}

func (m *ExchangeManager) initCollectorAndAggregator() {
	m.collector = NewCollector(m.cache, m.logger)
	m.aggregator = NewAggregator(m.store, m.cache, m.cfg.Aggregator, m.logger)
}

func (m *ExchangeManager) getWorkerPoolOutputs() []<-chan *domain.PriceData {
	var chans []<-chan *domain.PriceData
	for _, pool := range m.workerPools {
		chans = append(chans, pool.Output())
	}
	return chans
}
