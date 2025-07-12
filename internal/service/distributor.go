package service

import (
	"context"

	"marketflow/internal/domain"
)

type workerPool interface {
	Input() chan<- *domain.PriceData
}

// Distributor
type Distributor struct {
	workerPool workerPool
	in         <-chan *domain.PriceData
}

func NewDistriubtor(workerPool *WorkerPool, in <-chan *domain.PriceData) *Distributor {
	return &Distributor{
		workerPool: workerPool,
		in:         in,
	}
}

// FanOut reads data from incoming channel. And sends to pool of workers
func (d *Distributor) FanOut(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-d.in:
				if !ok {
					return // input channel closed
				}
				select {
				case d.workerPool.Input() <- data:
					// sent successfully
				case <-ctx.Done():
					return
				}
			}
		}
	}()
}
