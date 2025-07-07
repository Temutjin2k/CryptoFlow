package service

import (
	"context"
	"marketflow/internal/domain"
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"math/rand"
	"sync"
	"time"
)

type ExchangeManager struct {
	collector       ports.Collector
	exchangeSources []ports.ExchangeClient
	logger          logger.Logger
}

func NewExchangeManager(collector ports.Collector, exchange []ports.ExchangeClient, logger logger.Logger) *ExchangeManager {
	return &ExchangeManager{
		collector:       collector,
		exchangeSources: exchange,
		logger:          logger,
	}
}

func (e *ExchangeManager) Start(ctx context.Context) error {
	var allExchanges []<-chan domain.PriceData
	for _, source := range e.exchangeSources {
		prices, _ := source.Start(ctx)
		allExchanges = append(allExchanges, prices)
	}

	var fanOutChannels [][]<-chan domain.PriceData
	for _, sourceChan := range allExchanges {
		fanOutChannels = append(fanOutChannels, FanOut(sourceChan, 5)) // 5 workers на источник
	}

	// Запуск worker'ов
	var workerResults []<-chan domain.PriceData
	for _, foChans := range fanOutChannels {
		workerResults = append(workerResults, startWorkers(len(foChans), foChans))
	}

	// Fan-in результатов
	aggregatedResults := FanIn(workerResults...)

	e.collector.Start(aggregatedResults)

	return nil
}

// FanIn implements "FanIn" concurency pattern
func FanIn(channels ...<-chan domain.PriceData) <-chan domain.PriceData {
	var wg sync.WaitGroup
	out := make(chan domain.PriceData)

	output := func(c <-chan domain.PriceData) {
		defer wg.Done()
		for price := range c {
			out <- price
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FanOut implements "FanOut" concurency pattern
func FanOut(source <-chan domain.PriceData, numWorkers int) []<-chan domain.PriceData {
	outputs := make([]<-chan domain.PriceData, numWorkers)

	for i := range numWorkers {
		output := make(chan domain.PriceData)
		outputs[i] = output

		go func(out chan<- domain.PriceData) {
			defer close(out)
			for price := range source {
				out <- price
			}
		}(output)
	}

	return outputs
}

func worker(prices <-chan domain.PriceData, results chan<- domain.PriceData) {
	for price := range prices {
		// Имитация обработки
		// processed := domain.PriceData{
		// 	Price:    price,
		// 	IsValid:  validatePrice(price),
		// 	Metadata: enrichData(price),
		// }

		// Добавляем задержку для имитации работы
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

		results <- price
	}
}

func startWorkers(numWorkers int, inputs []<-chan domain.PriceData) <-chan domain.PriceData {
	results := make(chan domain.PriceData)
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(workerID int, in <-chan domain.PriceData) {
			defer wg.Done()
			worker(in, results)
		}(i+1, input)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
