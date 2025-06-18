package exchange

import (
	"marketflow/internal/domain"
	"math/rand"
	"time"
)

func GenerateTestData() <-chan domain.PriceData {
	out := make(chan domain.PriceData)

	symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT", "TONUSDT", "DOGEUSDT"}

	go func() {
		defer close(out)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		for {
			for _, symbol := range symbols {
				price := generateRandomPrice(symbol, r)
				data := domain.PriceData{
					Exchange:  "test",
					Symbol:    symbol,
					Price:     price,
					Timestamp: time.Now(),
				}

				out <- data
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return out
}

func generateRandomPrice(symbol string, r *rand.Rand) float64 {
	switch symbol {
	case "BTCUSDT":
		return 105000 + r.Float64()*500
	case "ETHUSDT":
		return 4200 + r.Float64()*50
	case "SOLUSDT":
		return 100 + r.Float64()*5
	case "TONUSDT":
		return 4 + r.Float64()*0.05
	case "DOGEUSDT":
		return 0.27 + r.Float64()*0.01
	default:
		return 100 + r.Float64()*10
	}
}
