package errors

import (
	"errors"
	"net/http"

	"github.com/spearszx/go-wallet/monolith/internal/apperrors"
)

const (
	unknownInternalFailureMessage = "internal server failure"
)

type APIError struct {
	Status  int
	Message string
}

func (e APIError) Error() string { return e.Message }

func InternalFailure() APIError {
	return APIError{
		Status:  http.StatusInternalServerError,
		Message: unknownInternalFailureMessage,
	}
}

func statusForCategory(c apperrors.Category) int {
	switch c {
	case apperrors.CategoryBadRequest:
		return http.StatusBadRequest
	case apperrors.CategoryInternal:
		return http.StatusInternalServerError
	case apperrors.CategoryNotFound:
		return http.StatusNotFound
	case apperrors.CategoryConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func FromError(e error) APIError {
	var cat apperrors.Categorized
	if !errors.As(e, &cat) {
		return InternalFailure()
	}

	return APIError{
		Status:  statusForCategory(cat.Category()),
		Message: cat.Error(),
	}
}

func New(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}
