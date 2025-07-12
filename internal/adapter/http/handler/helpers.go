package handler

import (
	"encoding/json"
	"errors"
	"maps"
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

func notFoundErrorResponse(w http.ResponseWriter) {
	errorResponse(w, http.StatusNotFound, "requested resource not found")
}

func parsePeriod(period string) (time.Duration, string, error) {
	if period == "" {
		period = "1m"
	}

	parsed, err := time.ParseDuration(period)
	if err != nil {
		return -1, "", err
	}

	if parsed <= 0 {
		return -1, "", errors.New("invalid period format. should be positive non-zero value (e.g. 1s, 5s, 1m, 3m)")
	}

	if parsed > time.Minute*5 {
		return -1, "", errors.New("must be less than 5 minutes")
	}

	return parsed, period, nil
}
