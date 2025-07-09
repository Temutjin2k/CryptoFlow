package service

import (
	"context"
	"fmt"
	"marketflow/internal/domain"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
)

type ExchangeManager struct {
	exchangeSources []ports.ExchangeClient
	distributors    []ports.Distributor
	workerPools     []ports.WorkerPool
	aggregator      ports.Aggregator
	collector       ports.Collector

	logger logger.Logger
}

func NewExchangeManager(
	exchange []ports.ExchangeClient,
	collector ports.Collector,
	aggregator ports.Aggregator,
	logger logger.Logger,
) *ExchangeManager {
	return &ExchangeManager{
		exchangeSources: exchange,
		aggregator:      aggregator,
		collector:       collector,
		logger:          logger,
	}
}

// Start
func (m *ExchangeManager) Start(ctx context.Context) error {
	// Iterating accross exchanges and starting.
	for _, source := range m.exchangeSources {
		pricesCh, err := source.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start source: %w", err)
		}
		workerPool := NewWorkerPool(10, m.logger)
		distributor := NewDistriubtor(workerPool, pricesCh)

		// starting worker pools.
		distributor.FanOut(ctx)

		// Saving distributors
		m.distributors = append(m.distributors, distributor)
	}

	resultch := m.aggregator.FanIn(m.getWorkerPoolOutputs()...)

	m.collector.Start(resultch)

	return nil
}

func (m *ExchangeManager) Close() error {
	return nil
}

func (m *ExchangeManager) getWorkerPoolOutputs() []<-chan *domain.PriceData {
	channels := make([]<-chan *domain.PriceData, len(m.workerPools))
	for _, v := range m.workerPools {
		channels = append(channels, v.Output())
	}

	return channels
}
