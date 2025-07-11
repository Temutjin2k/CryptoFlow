package handler

import (
	"encoding/json"
	"errors"
	"maps"
	"marketflow/pkg/logger"
	"net/http"
	"time"
)

type envelope map[string]any

func errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}

	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	err := writeJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return errors.New("failed to encode json")
	}

	js = append(js, '\n')

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func internalErrorResponse(w http.ResponseWriter, message any) {
	errorResponse(w, http.StatusInternalServerError, message)
}

func parsePeriod(w http.ResponseWriter, r *http.Request, log logger.Logger) (time.Duration, bool) {
	period := r.URL.Query().Get("period")
	parsed, err := time.ParseDuration(period)
	if err != nil {
		log.Error(r.Context(), "invalid period format", "error", err)
		errorResponse(w, http.StatusBadRequest, "invalid period format")
		return 0, false
	}
	return parsed, true
}
