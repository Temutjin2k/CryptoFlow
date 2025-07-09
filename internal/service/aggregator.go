package service

import (
	"fmt"
	"marketflow/internal/domain"
	"sync"
)

type Aggregator struct {
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) FanIn(inputs ...<-chan *domain.PriceData) <-chan *domain.PriceData {
	output := make(chan *domain.PriceData)

	var wg sync.WaitGroup
	for i, input := range inputs {
		wg.Add(1)
		go func() {
			for value := range input {
				output <- value
			}
			fmt.Printf("done merging source %d\n", i)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(output) // this is important!!!
	}()

	return output
}
