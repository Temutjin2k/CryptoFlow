package service

import (
	"marketflow/internal/domain"
	"sync"
)

// Aggregator represent structure to calculate avg, min, max prices.
type Aggregator struct {
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) FanIn(inputs ...<-chan *domain.PriceData) <-chan *domain.PriceData {
	output := make(chan *domain.PriceData)

	var wg sync.WaitGroup
	for _, input := range inputs {
		wg.Add(1)
		go func() {
			for value := range input {
				output <- value
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
