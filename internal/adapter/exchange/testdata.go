package exchange

import (
	"context"
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"math/rand"
	"time"
)

type TestExchangeSource struct {
	name string
	stop chan struct{}
}

func NewTestExchange(name string) *TestExchangeSource {
	return &TestExchangeSource{
		name: name,
		stop: make(chan struct{}),
	}
}

func (t *TestExchangeSource) Name() string {
	return t.name
}

func (t *TestExchangeSource) Start(ctx context.Context) (<-chan *domain.PriceData, error) {
	out := make(chan *domain.PriceData)

	go func() {
		symbols := []types.Symbol{
			types.BTCUSDT,
			types.ETHUSDT,
			types.SOLUSDT,
			types.TONUSDT,
			types.DOGEUSDT,
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for {
			select {
			case <-ctx.Done():
				close(out)
				return
			case <-t.stop:
				close(out)
				return
			default:
				for _, symbol := range symbols {
					out <- &domain.PriceData{
						Exchange:  types.TestExchange,
						Symbol:    symbol,
						Price:     generateRandomPrice(symbol, r),
						Timestamp: time.Now(),
					}
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return out, nil
}

func (t *TestExchangeSource) Close() error {
	close(t.stop)
	return nil
}

// generateRandomPrice generates random price fluctuation (±15%)
func generateRandomPrice(symbol types.Symbol, r *rand.Rand) float64 {
	switch symbol {
	case types.BTCUSDT:
		return 10 + r.Float64()*2 // 10 - 12
	case types.ETHUSDT:
		return 5 + r.Float64()*2 // 5 - 7
	case types.SOLUSDT:
		return 2 + r.Float64()*1.5 // 2 - 3.5
	case types.TONUSDT:
		return 1 + r.Float64()*1 // 1 - 2
	case types.DOGEUSDT:
		return 0.1 + r.Float64()*0.1 // 0.1 - 0.2
	default:
		return 5 + r.Float64()*1 // 5 - 6
	}
}
