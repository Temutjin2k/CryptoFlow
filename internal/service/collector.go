package service

import (
	"context"
	"marketflow/internal/domain"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"time"
)

type Collector struct {
	cache  ports.Cache
	store  ports.MarketRepository
	logger logger.Logger
}

func NewCollector(cache ports.Cache, store ports.MarketRepository, logger logger.Logger) *Collector {
	return &Collector{
		cache:  cache,
		store:  store,
		logger: logger,
	}
}

func (c *Collector) Start(processedPrices <-chan domain.PriceData) {
	const fn = "collector.Start"
	log := c.logger.GetSlogLogger().With("fn", fn)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var count int
	var lastPrice domain.PriceData
	go func() {
		for {
			select {
			case price, ok := <-processedPrices:
				if !ok {
					log.Error("processing completed", "total_prices", count)
					return
				}

				count++
				lastPrice = price

				// TODO set price data correctly
				if err := c.cache.SetLatest(context.Background(), price, time.Minute); err != nil {
					log.Error("cache store failed", "error", err)
				}

			case <-ticker.C:
				log.Error("processing status",
					"total_processed", count,
					"last_price", lastPrice)
			}
		}
	}()

}
