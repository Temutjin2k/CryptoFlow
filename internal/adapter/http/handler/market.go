package handler

import (
	"marketflow/pkg/logger"
	"net/http"
)

type Market struct {
	log logger.Logger
}

func NewMarket(log logger.Logger) *Market {
	return &Market{
		log: log,
	}
}

// LATEST

// GET localhost:8080/prices/latest/{symbol}
func (h *Market) LatestPrice(w http.ResponseWriter, r *http.Request) { //latest price over a period

}

// GET localhost:8080/prices/latest/{exchange}/{symbol}
func (h *Market) LatestPriceByExchange(w http.ResponseWriter, r *http.Request) { //latest price over a period by exchange

}

// HIGHEST

// GET localhost:8080/prices/highest/{symbol}
func (h *Market) HighestPrice(w http.ResponseWriter, r *http.Request) { //highest price over a period

}

// GET localhost:8080/prices/highest/{exchange}/{symbol}
func (h *Market) HighestPriceByExchange(w http.ResponseWriter, r *http.Request) { //highest price over a period by exchange

}

// GET localhost:8080/prices/highest/{symbol}?period={duration}
func (h *Market) HighestPriceDuration(w http.ResponseWriter, r *http.Request) { //highest price within a duration

}

// GET localhost:8080/prices/highest/{exchange}/{symbol}?period={duration}
func (h *Market) HighestPriceDurationByExchange(w http.ResponseWriter, r *http.Request) { //highest price within a duration by exchange

}

// LOWEST

// GET localhost:8080/prices/lowest/{symbol}
func (h *Market) LowestPrice(w http.ResponseWriter, r *http.Request) { //lowest price over a period

}

// GET localhost:8080/prices/lowest/{exchange}/{symbol}
func (h *Market) LowestPriceByExchange(w http.ResponseWriter, r *http.Request) { //lowest price by exchange

}

// GET localhost:8080/prices/lowest/{symbol}?period={duration}
func (h *Market) LowestPriceDuration(w http.ResponseWriter, r *http.Request) { //lowest price within a duration

}

// GET localhost:8080/prices/lowest/{exchange}/{symbol}?period={duration}
func (h *Market) LowestPriceDurationByExchange(w http.ResponseWriter, r *http.Request) { //lowest price within a duration by exchange

}

// AVERAGE

// GET localhost:8080/prices/average/{symbol}
func (h *Market) AveragePrice(w http.ResponseWriter, r *http.Request) { //average price over a period

}

// GET localhost:8080/prices/average/{exchange}/{symbol}?period={duration}
func (h *Market) AveragePriceDurationByExchange(w http.ResponseWriter, r *http.Request) { //average price within a duration by exchange

}

// GET localhost:8080/prices/average/{exchange}/{symbol}
func (h *Market) AveragePriceByExchange(w http.ResponseWriter, r *http.Request) { //average price by exchange

}
