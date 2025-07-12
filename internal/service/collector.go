package service

import (
	"context"
	"errors"
	"time"

	"marketflow/internal/domain"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
)

type Collector struct {
	cache ports.Cache

	cancelFunc context.CancelFunc
	doneChan   chan struct{}

	logger logger.Logger
}

func NewCollector(cache ports.Cache, logger logger.Logger) *Collector {
	return &Collector{
		cache:    cache,
		doneChan: make(chan struct{}),
		logger:   logger,
	}
}

// Start starts to listens incoming channel
func (c *Collector) Start(ctx context.Context, processedPrices <-chan *domain.PriceData) {
	ctx, cancel := context.WithCancel(ctx)
	c.cancelFunc = cancel

	go c.run(ctx, processedPrices)
}

func (c *Collector) run(ctx context.Context, processedPrices <-chan *domain.PriceData) {
	defer close(c.doneChan)

	const fn = "collector.run"
	log := c.logger.GetSlogLogger().With("fn", fn)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var count int
	// var lastPrice *domain.PriceData

	for {
		select {
		case <-ctx.Done():
			log.Info("shutting down collector")
			return

		case price, ok := <-processedPrices:
			if !ok {
				log.Info("collector's input channel closed", "total_prices", count)
				return
			}

			count++
			// lastPrice = price

			if err := c.cache.SetLatest(ctx, price, time.Minute); err != nil {
				log.Error("cache store failed", "error", err)
			}

			if err := c.cache.StoreHistory(ctx, price); err != nil {
				log.Error("history store failed", "error", err)
			}

		case <-ticker.C:
			// log.Info("processing status",
			// 	"total_processed", count,
			// 	"last_price", lastPrice)
		}
	}
}

// Cancel gracefully shutdowns collector
func (c *Collector) Cancel() error {
	if c.cancelFunc != nil {
		c.cancelFunc()

		// Waiting confirmation to close collector
		select {
		case <-c.doneChan:
			return nil
		case <-time.After(5 * time.Second):
			return errors.New("timeout waiting for collector to stop")
		}
	}
	return nil
}
