package service

import (
	"context"
	"marketflow/internal/domain"
	"time"
)

type Aggregator struct {
	service *Market
}

func NewAggregator(service *Market) *Aggregator {
	return &Aggregator{
		service: service,
	}
}

func (a *Aggregator) FanIn(ctx context.Context, _ ...<-chan *domain.PriceData) {
	if a.service == nil {
		panic("aggregator service is nil")
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.service.AggregateAndStore(ctx)
			}
		}
	}()
}
