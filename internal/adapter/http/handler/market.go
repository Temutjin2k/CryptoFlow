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

func (h *Market) LatestPrice(w http.ResponseWriter, r *http.Request) {

}

func (h *Market) LatestPriceByExchange(w http.ResponseWriter, r *http.Request) {

}

func (h *Market) HighestPrice(w http.ResponseWriter, r *http.Request) {

}

func (h *Market) HighestPriceByExchange(w http.ResponseWriter, r *http.Request) {

}

func (h *Market) LowestPrice(w http.ResponseWriter, r *http.Request) {

}

func (h *Market) LowestPriceByExchange(w http.ResponseWriter, r *http.Request) {

}

func (h *Market) AveragePrice(w http.ResponseWriter, r *http.Request) {

}

func (h *Market) AveragePriceByExchange(w http.ResponseWriter, r *http.Request) {

}
