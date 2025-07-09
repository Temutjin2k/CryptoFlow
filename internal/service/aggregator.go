package service

import (
	"marketflow/internal/domain"
	"sync"
)

// Aggregator represent structure to calculate avg, min, max prices.
type Aggregator struct {
	input chan *domain.PriceData
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) FanIn(inputs ...<-chan *domain.PriceData) {
	var wg sync.WaitGroup
	for _, input := range inputs {
		wg.Add(1)
		go func() {
			for value := range input {
				a.input <- value
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(a.input)
	}()
}

func (a *Aggregator) listenAndServe() {

}

func (a *Aggregator) Input() <-chan *domain.PriceData {
	return a.input
}
