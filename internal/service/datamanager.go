package service

import (
	"context"
	"fmt"
	"marketflow/internal/domain"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"sync"
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

func (m *ExchangeManager) Start(ctx context.Context) error {
	for _, source := range m.exchangeSources {
		m.logger.Info(ctx, "starting source", "name", source.Name())
		pricesCh, err := source.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start source: %w", err)
		}

		workerPool := NewWorkerPool(source.Name(), m.workerCount, m.logger)
		m.workerPools = append(m.workerPools, workerPool)

		distributor := NewDistriubtor(workerPool, pricesCh)
		m.distributors = append(m.distributors, distributor)

		workerPool.Start(ctx)
		distributor.FanOut(ctx)
	}

	m.aggregator.FanIn(ctx, m.getWorkerPoolOutputs()...)

	merged := mergeChannels(ctx, m.getWorkerPoolOutputs()...)
	m.collector.Start(ctx, merged)

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

func (m *ExchangeManager) getWorkerPoolOutputs() []<-chan *domain.PriceData {
	var chans []<-chan *domain.PriceData
	for _, pool := range m.workerPools {
		chans = append(chans, pool.Output())
	}
	return chans
}

func mergeChannels(ctx context.Context, chans ...<-chan *domain.PriceData) <-chan *domain.PriceData {
	out := make(chan *domain.PriceData)
	var wg sync.WaitGroup

	for _, ch := range chans {
		wg.Add(1)
		go func(c <-chan *domain.PriceData) {
			defer wg.Done()
			for val := range c {
				select {
				case out <- val:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
