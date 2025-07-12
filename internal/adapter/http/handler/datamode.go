package handler

import (
	"marketflow/internal/ports"
	"marketflow/pkg/logger"
	"net/http"
)

type DataMode struct {
	manager     ports.ExchangeManager
	testSources []ports.ExchangeSource
	liveSources []ports.ExchangeSource
	log         logger.Logger
}

func NewDataMode(manager ports.ExchangeManager, test, live []ports.ExchangeSource, log logger.Logger) *DataMode {
	return &DataMode{
		manager:     manager,
		testSources: test,
		liveSources: live,
		log:         log,
	}
}

func (h *DataMode) TestMode(w http.ResponseWriter, r *http.Request) {
	err := h.manager.SwitchToTest(h.testSources)
	if err != nil {
		h.log.Error(r.Context(), "failed to switch to test mode", "error", err)
		internalErrorResponse(w, "failed to switch to test mode")
		return
	}
	writeJSON(w, http.StatusOK, envelope{"message": "switched to test mode"}, nil)
}

func (h *DataMode) LiveMode(w http.ResponseWriter, r *http.Request) {
	err := h.manager.SwitchToLive(h.liveSources)
	if err != nil {
		h.log.Error(r.Context(), "failed to switch to live mode", "error", err)
		internalErrorResponse(w, "failed to switch to live mode")
		return
	}
	writeJSON(w, http.StatusOK, envelope{"message": "switched to live mode"}, nil)
}
