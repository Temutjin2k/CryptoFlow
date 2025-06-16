package model

type PriceData struct {
	Symbol    string  `json:"symbol"`
	Exchange  string  `json:"exchange"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}
