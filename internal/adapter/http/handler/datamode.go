package handler

import (
	"errors"
	"net/http"

	"marketflow/internal/domain"
	"marketflow/pkg/logger"
)

type ModeSwitcher interface {
	SwitchToTest() error
	SwitchToLive() error
}

type DataMode struct {
	mode ModeSwitcher
	log  logger.Logger
}

func NewDataMode(manager ModeSwitcher, log logger.Logger) *DataMode {
	return &DataMode{
		mode: manager,
		log:  log,
	}
}

func (h *DataMode) TestMode(w http.ResponseWriter, r *http.Request) {
	if err := h.mode.SwitchToTest(); err != nil {
		if errors.Is(err, domain.ErrAlreadyOnTestMode) {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		h.log.Error(r.Context(), "failed to switch to test mode", "error", err)
		internalErrorResponse(w, "failed to switch to test mode")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"message": "switched to test mode"}, nil)
}

func (h *DataMode) LiveMode(w http.ResponseWriter, r *http.Request) {
	if err := h.mode.SwitchToLive(); err != nil {
		if errors.Is(err, domain.ErrAlreadyOnLiveMode) {
			errorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		h.log.Error(r.Context(), "failed to switch to live mode", "error", err)
		internalErrorResponse(w, "failed to switch to live mode")
		return
	}

	writeJSON(w, http.StatusOK, envelope{"message": "switched to live mode"}, nil)
}
