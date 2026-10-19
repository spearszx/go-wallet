package service

import "github.com/spearszx/go-wallet/monolith/internal/apperrors"

var (
	ErrUserNofFound = apperrors.New(apperrors.CategoryNotFound, "user not found")
	ErrInvalidInput = apperrors.New(apperrors.CategoryBadRequest, "invalid input")
)
