package middleware

import (
	"net/http"

	"github.com/spearszx/go-wallet/monolith/internal/api"
	"github.com/spearszx/go-wallet/monolith/internal/errors"
)

const methodNotAllowedMessage = "method is not allowed"

func MethodPost(fn api.HttpFunc) api.HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			return errors.New(http.StatusMethodNotAllowed, methodNotAllowedMessage)
		}

		return fn(w, r)
	}
}

func MethodDelete(fn api.HttpFunc) api.HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodDelete {
			return errors.New(http.StatusMethodNotAllowed, methodNotAllowedMessage)
		}

		return fn(w, r)
	}
}

func MethodPatch(fn api.HttpFunc) api.HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPatch {
			return errors.New(http.StatusMethodNotAllowed, methodNotAllowedMessage)
		}

		return fn(w, r)
	}
}

func MethodGet(fn api.HttpFunc) api.HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return errors.New(http.StatusMethodNotAllowed, methodNotAllowedMessage)
		}

		return fn(w, r)
	}
}

func MethodPut(fn api.HttpFunc) api.HttpFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPut {
			return errors.New(http.StatusMethodNotAllowed, methodNotAllowedMessage)
		}

		return fn(w, r)
	}
}
