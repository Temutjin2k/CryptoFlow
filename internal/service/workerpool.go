package service

import (
	"context"
	"errors"
	"marketflow/internal/domain"
	"marketflow/pkg/logger"
	"sync"
)

// WorkerPool represents a pool of workers for processing PriceData
type WorkerPool struct {
	name string

	workerCount int
	inputChan   chan *domain.PriceData
	outputChan  chan *domain.PriceData
	wg          sync.WaitGroup

	log logger.Logger
}

func NewWorkerPool(name string, workerCount int, log logger.Logger) *WorkerPool {
	return &WorkerPool{
		name:        name,
		workerCount: workerCount,
		inputChan:   make(chan *domain.PriceData, 100),
		outputChan:  make(chan *domain.PriceData, 100),
		log:         log,
	}
}

// Start launches all workers in the pool
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx)
	}

	wp.log.Info(ctx, "started workers", "pool_name", wp.name, "count", wp.workerCount)

	go func() {
		wp.wg.Wait()
		close(wp.outputChan)
	}()
}

// worker processes incoming data
func (wp *WorkerPool) worker(ctx context.Context) {
	defer wp.wg.Done()
	log := wp.log.GetSlogLogger().With("pool_name", wp.name)
	log.Info("worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info("worker stopped by context")
			return
		case priceData, ok := <-wp.inputChan:
			if !ok {
				return
			}
			processed, err := wp.processPriceData(priceData)
			if err != nil {
				log.Error("failed to process price data")
				continue
			}

			wp.outputChan <- processed
		}
	}
}

// processPriceData contains the logic for processing/validating PriceData
func (wp *WorkerPool) processPriceData(data *domain.PriceData) (*domain.PriceData, error) {
	ok, err := data.IsValid()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed")
	}

	// Return the processed data (unchanged in this example)
	return data, nil
}

// Input returns a channel to send data into the pool
func (wp *WorkerPool) Input() chan<- *domain.PriceData {
	return wp.inputChan
}

// Output returns a channel to receive processed data from the pool
func (wp *WorkerPool) Output() <-chan *domain.PriceData {
	return wp.outputChan
}

// Close closes the input channel (signals workers to stop)
func (wp *WorkerPool) Close() {
	wp.log.Info(context.Background(), "closing worker pool")

	// Check if channel is already closed
	select {
	case _, ok := <-wp.inputChan:
		if !ok {
			return
		}
	default:
	}

	close(wp.inputChan)
}
