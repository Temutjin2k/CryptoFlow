package service

import (
	"context"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"sync"
	"time"
)

type Aggregator struct {
	storage        ports.MarketRepository
	cache          ports.Cache
	tickerDuration time.Duration

	logger logger.Logger
}

func NewAggregator(storage ports.MarketRepository, cache ports.Cache, tickerDuration time.Duration, logger logger.Logger) *Aggregator {
	return &Aggregator{
		storage:        storage,
		cache:          cache,
		tickerDuration: tickerDuration,

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
		ticker := time.NewTicker(a.tickerDuration)
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

func aggregate(values []float64) (min, max, avg float64) {
	min, max = values[0], values[0]
	sum := 0.0

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}

	avg = sum / float64(len(values))
	return
}
