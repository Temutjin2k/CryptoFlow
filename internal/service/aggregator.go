package service

import (
	"context"
	"marketflow/config"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"sync"
	"time"
)

type Aggregator struct {
	storage ports.MarketRepository
	cache   ports.Cache

	cfg    config.Aggregator
	logger logger.Logger
}

func NewAggregator(storage ports.MarketRepository, cache ports.Cache, cfg config.Aggregator, logger logger.Logger) *Aggregator {
	return &Aggregator{
		storage: storage,
		cache:   cache,
		cfg:     cfg,

		logger: logger,
	}
}

func (a *Aggregator) FanIn(ctx context.Context, inputs ...<-chan *domain.PriceData) <-chan *domain.PriceData {
	output := make(chan *domain.PriceData)

	var wg sync.WaitGroup

	for _, ch := range inputs {
		wg.Add(1)

		go func(in <-chan *domain.PriceData) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func (a *Aggregator) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(a.cfg.TickerDuration)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.aggregateAndStore(ctx)
			}
		}
	}()
}

// aggregateAndStore gets aggregated data and sends it to database
func (s *Aggregator) aggregateAndStore(ctx context.Context) {
	s.logger.Info(ctx, "Running AggregateAndStore")
	exchanges := types.ValidExchanges
	symbols := types.ValidSymbols

	stats := []*domain.PriceStats{}

	for _, exchange := range exchanges {
		for _, symbol := range symbols {
			values, err := s.cache.GetPriceInPeriod(ctx, exchange, symbol, time.Minute)
			if len(values) == 0 {
				s.logger.Warn(ctx, "No prices found wait for 1 minute", "exchange", exchange, "symbol", symbol, "len", len(values), "err", err)
				continue
			}
			if err != nil {
				s.logger.Error(ctx, "failed to get prices from cache", "exchange", exchange, "symbol", symbol, "error", err)
				continue
			}

			min, max, avg := aggregate(values)

			stat := &domain.PriceStats{
				Exchange:  exchange,
				Pair:      symbol,
				Timestamp: time.Now(),
				Average:   avg,
				Min:       min,
				Max:       max,
			}

			stats = append(stats, stat)
		}
	}

	// Saving to the database
	if err := s.storage.StoreStats(ctx, stats); err != nil {
		s.logger.Error(ctx, "failed to save stats")
	}
}
