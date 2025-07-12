package server

import (
	"net/http"
	"time"
)

// LoggingMiddleware wraps an http.Handler and logs requests
func (m *API) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture the status code
		lrw := NewLoggingResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log the request details
		m.log.Info(r.Context(), "request details",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lrw.statusCode,
			"duration", time.Since(start),
			"remote_ip", r.RemoteAddr,
		)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK} // Default to 200 OK
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
