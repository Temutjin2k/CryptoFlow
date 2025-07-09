package service

import (
	"context"
	"fmt"
	"marketflow/internal/domain"
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

	workerCount int // number of workers per each distributor
	logger      logger.Logger
}

func NewExchangeManager(
	exchange []ports.ExchangeSource,
	collector ports.Collector,
	aggregator ports.Aggregator,
	workerCount int,

	logger logger.Logger,
) *ExchangeManager {
	return &ExchangeManager{
		exchangeSources: exchange,
		aggregator:      aggregator,
		collector:       collector,
		workerCount:     workerCount,
		logger:          logger,
	}
}

// Start
func (m *ExchangeManager) Start(ctx context.Context) error {
	// Iterating accross exchanges and starting sources.
	for _, source := range m.exchangeSources {
		pricesCh, err := source.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start source: %w", err)
		}

		workerPool := NewWorkerPool(source.Name(), m.workerCount, m.logger)
		m.workerPools = append(m.workerPools, workerPool)

		distributor := NewDistriubtor(workerPool, pricesCh)
		m.distributors = append(m.distributors, distributor)

		// starting worker pools.
		workerPool.Start(ctx)
		//starting distributor.
		distributor.FanOut(ctx)
	}

	m.aggregator.FanIn(m.getWorkerPoolOutputs()...)

	resultch := m.aggregator.Input()

	m.collector.Start(ctx, resultch)

	return nil
}

func (m *ExchangeManager) Close() error {
	const fn = "ExchangeManager.Close"
	log := m.logger.GetSlogLogger().With("fn", fn)

	// Closing the exchange sources.
	for _, source := range m.exchangeSources {
		if err := source.Close(); err != nil {
			log.Warn("failed to close exchange", "error", err)
		}
	}

	// TODO: maybe close aggregator, but not sure.

	// Closing collector
	if err := m.collector.Cancel(); err != nil {
		log.Warn("failed to cancel collector", "error", err)
	}

	// Closing the workerPools
	for _, pool := range m.workerPools {
		pool.Close()
	}

	return nil
}

func (m *ExchangeManager) getWorkerPoolOutputs() []<-chan *domain.PriceData {
	var chans []<-chan *domain.PriceData
	for _, pool := range m.workerPools {
		chans = append(chans, pool.Output())
	}

	return chans
}
