package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWritter struct {
	http.ResponseWriter
	Status int
}

func (rw *responseWritter) WriteHeader(code int) {
	rw.Status = code
	rw.WriteHeader(code)
}

func Logger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWritter{ResponseWriter: w, Status: http.StatusOK}
		next.ServeHTTP(rw, r)

		slog.Info("http",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.Status),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
			slog.String("request_id", IDFromContext(r.Context())),
		)
	})
}
