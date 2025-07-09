package service

import (
	"context"
	"marketflow/internal/domain"
)

type Distributor struct {
	workerPool *WorkerPool
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
		defer d.workerPool.Close()
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

func (d *Distributor) Close() error {
	return nil
}
