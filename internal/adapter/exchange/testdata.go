package exchange

import (
	"marketflow/internal/domain"
	"marketflow/internal/domain/types"
	"math/rand"
	"time"
)

type testExchange struct {
}

// GenerateTestData returns channel with generated data
func GenerateTestData() <-chan *domain.PriceData {
	out := make(chan *domain.PriceData)

	symbols := []types.Symbol{
		types.BTCUSDT,
		types.ETHUSDT,
		types.SOLUSDT,
		types.TONUSDT,
		types.DOGEUSDT,
	}

	go func() {
		defer close(out)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for {
			for _, symbol := range symbols {
				price := generateRandomPrice(symbol, r)
				data := &domain.PriceData{
					Exchange:  "test",
					Symbol:    symbol,
					Price:     price,
					Timestamp: time.Now(),
				}

				out <- data
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return out
}

// generateRandomPrice generates random price fluctuation (±15%)
func generateRandomPrice(symbol types.Symbol, r *rand.Rand) float64 {
	switch symbol {
	case types.BTCUSDT:
		return 105000 * (1 + (r.Float64()-0.5)*0.3)
	case types.ETHUSDT:
		return 4200 * (1 + (r.Float64()-0.5)*0.3)
	case types.SOLUSDT:
		return 100 * (1 + (r.Float64()-0.5)*0.3)
	case types.TONUSDT:
		return 4 * (1 + (r.Float64()-0.5)*0.3)
	case types.DOGEUSDT:
		return 0.27 * (1 + (r.Float64()-0.5)*0.3)
	default:
		return 100 * (1 + (r.Float64()-0.5)*0.3)
	}
}
