package middleware

import (
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// slog.Error("panic recovered",
				// 	"panic", rec,
				// 	"method", r.Method,
				// 	"path", r.URL.Path,
				// 	"stack", string(debug.Stack()),
				// )
			}

		}()

		next.ServeHTTP(w, r)
	})
}
