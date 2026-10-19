package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spearszx/go-wallet/monolith/internal/api"
	httpErrors "github.com/spearszx/go-wallet/monolith/internal/errors"
)

func ErrorHandler(fn api.HttpFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var apiErr httpErrors.APIError
			if errors.As(err, &apiErr) {
				writeError(w, apiErr.Status, apiErr.Message)
			}

			writeError(w, http.StatusInternalServerError, httpErrors.InternalFailure().Message)
		}
	})
}

func writeError(w http.ResponseWriter, status int, message string) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(map[string]any{"success": false, "message": message})
}
