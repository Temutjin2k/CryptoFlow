package exchange

import (
	"context"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"math/rand"
	"time"
)

type TestExchangeSource struct {
	name   types.Exchange
	cancel context.CancelFunc
}

func NewTestExchange(name types.Exchange) *TestExchangeSource {
	return &TestExchangeSource{
		name: name,
	}
}

func (t *TestExchangeSource) Start(ctx context.Context) (<-chan *domain.PriceData, error) {
	out := make(chan *domain.PriceData)

	ctx, cancel := context.WithCancel(ctx)
	t.cancel = cancel

	go func() {
		defer close(out) // Ensure channel is always closed
		symbols := types.ValidSymbols

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return // Context cancellation
			case <-ticker.C:
				for _, symbol := range symbols {
					select {
					case out <- &domain.PriceData{
						Exchange:  t.name,
						Symbol:    symbol,
						Price:     generateRandomPrice(symbol, r),
						Timestamp: time.Now(),
					}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return out, nil
}

func (t *TestExchangeSource) Close() error {
	if t.cancel != nil {
		t.cancel() // canceling context to stop the gouroutine
	}
	return nil
}

func (t *TestExchangeSource) Name() string {
	return string(t.name)
}

// generateRandomPrice generates random price fluctuation (±15%)
func generateRandomPrice(symbol types.Symbol, r *rand.Rand) float64 {
	switch symbol {
	case types.BTCUSDT:
		return 100000 * (1 + (r.Float64()-0.5)*0.3)
	case types.ETHUSDT:
		return 5000 * (1 + (r.Float64()-0.5)*0.3)
	case types.SOLUSDT:
		return 200 * (1 + (r.Float64()-0.5)*0.3)
	case types.TONUSDT:
		return 99 * (1 + (r.Float64()-0.5)*0.3)
	case types.DOGEUSDT:
		return 0.27 * (1 + (r.Float64()-0.5)*0.3)
	default:
		return 66 * (1 + (r.Float64()-0.5)*0.3)
	}

}
