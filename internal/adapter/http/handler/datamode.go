package handler

import (
	"marketflow/pkg/logger"
	"net/http"
)

type DataMode struct {
	log logger.Logger
}

func NewDataMode(log logger.Logger) *DataMode {
	return &DataMode{
		log: log,
	}
}

func (h *DataMode) TestMode(w http.ResponseWriter, r *http.Request) {

}

func (h *DataMode) LiveMode(w http.ResponseWriter, r *http.Request) {

}
