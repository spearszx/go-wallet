package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/spearszx/go-wallet/monolith/internal/apperrors"
	"github.com/spearszx/go-wallet/monolith/internal/wallet/model"
	"github.com/spearszx/go-wallet/monolith/internal/wallet/repository"
)

type WalletService interface {
	Create(ctx context.Context, request *model.CreateWalletRequest) (*model.Wallet, error)
	FindByID(ctx context.Context, id string) (*model.Wallet, error)
	FindByUserID(ctx context.Context, userID string) (*model.Wallet, error)
	Update(ctx context.Context, id string, request *model.UpdatedWalletRequest) error
	Delete(ctx context.Context, id string) error
}

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) WalletService {
	return &walletService{
		walletRepo: walletRepo,
	}
}

func (s *walletService) Create(ctx context.Context, request *model.CreateWalletRequest) (*model.Wallet, error) {
	const op = "walletService.Create"

	wallet := model.Wallet{
		ID:       uuid.NewString(),
		UserID:   request.UserID,
		Balance:  request.Balance,
		Currency: request.Currency,
		Status:   request.Status,
		Version:  request.Version,
	}

	created, err := s.walletRepo.Create(ctx, &wallet)
	if err != nil {
		if errors.Is(err, repository.ErrWalletAlreadyExists) {
			return nil, apperrors.Wrap(apperrors.CategoryConflict, "wallet already exist", fmt.Errorf("%s: %w", op, err))
		}
		return nil, apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return created, nil
}

func newUserFromRequest(request *model.CreateWalletRequest) (*model.Wallet, error) {
	wallet := model.Wallet{
		ID:       uuid.NewString(),
		UserID:   request.UserID,
		Balance:  request.Balance,
		Currency: request.Currency,
		Status:   request.Status,
		Version:  request.Version,
	}

	return &wallet, nil
}

func (s *walletService) FindByID(ctx context.Context, id string) (*model.Wallet, error) {
	const op = "walletService.FindByID"

	user, err := s.walletRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			return nil, apperrors.Wrap(apperrors.CategoryNotFound, err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return nil, apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return user, nil
}

func (s *walletService) FindByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	const op = "walletService.FindByID"

	user, err := s.walletRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			return nil, apperrors.Wrap(apperrors.CategoryNotFound, err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return nil, apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return user, nil
}

func (s *walletService) Update(ctx context.Context, id string, request *model.UpdatedWalletRequest) error {
	const op = "walletService.Update"

	wallet, err := s.walletRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			return apperrors.Wrap(apperrors.CategoryNotFound,
				err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	updatWallet(request, wallet)

	if err = s.walletRepo.Update(ctx, wallet); err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			return apperrors.Wrap(apperrors.CategoryNotFound, err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return nil
}

func updatWallet(request *model.UpdatedWalletRequest, wallet *model.Wallet) {
	wallet.UserID = request.UserID
	wallet.Balance = request.Balance
	wallet.Currency = request.Currency
	wallet.Status = request.Status
	wallet.Version = request.Version
}

func (s *walletService) Delete(ctx context.Context, id string) error {
	const op = "walletService.Update"

	if err := s.walletRepo.Delete(ctx, id); err != nil {
		return apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return nil
}
