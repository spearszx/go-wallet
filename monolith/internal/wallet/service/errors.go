package service

import "github.com/spearszx/go-wallet/monolith/internal/apperrors"

var (
	ErrUserNofFound = apperrors.New(apperrors.CategoryNotFound, "wallet not found")
	ErrInvalidInput = apperrors.New(apperrors.CategoryBadRequest, "insufficient input")
)
