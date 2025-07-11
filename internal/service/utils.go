package service

import "marketflow/internal/domain"

func aggregate(values []*domain.PriceData) (min, max, avg float64) {
	min, max = values[0].Price, values[0].Price
	var sum float64

	for _, v := range values {
		if v.Price < min {
			min = v.Price
		}
		if v.Price > max {
			max = v.Price
		}
		sum += v.Price
	}

	avg = sum / float64(len(values))
	return
}

// aggregateAndPrice returns min, max, avg. WARNING values must from one exchange and one symbol.
func aggregateAndPrice(values []*domain.PriceData) (min, max, avg *domain.PriceData) {
	if len(values) == 0 {
		return nil, nil, nil
	}

	min = values[0]
	max = values[0]
	var sum float64

	for _, v := range values {
		if v.Price < min.Price {
			min = v
		}
		if v.Price > max.Price {
			max = v
		}
		sum += v.Price
	}

	avgPrice := sum / float64(len(values))

	avg = &domain.PriceData{
		Exchange:  values[len(values)-1].Exchange,
		Symbol:    values[len(values)-1].Symbol,
		Timestamp: values[len(values)-1].Timestamp,
		Price:     avgPrice,
	}

	return
}
