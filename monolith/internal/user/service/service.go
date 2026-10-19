package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/spearszx/go-wallet/monolith/internal/apperrors"
	"github.com/spearszx/go-wallet/monolith/internal/user/model"
	"github.com/spearszx/go-wallet/monolith/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, request *model.CreateUserRequest) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id string, request *model.UpdateUserRequest) error
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserSerivce(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, request *model.CreateUserRequest) (*model.User, error) {
	const op = "userService.Create"

	user, err := newUserFromRequest(request)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, apperrors.Wrap(apperrors.CategoryConflict, "user already exist", fmt.Errorf("%s: %w", op, err))
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return created, nil
}

func newUserFromRequest(request *model.CreateUserRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := model.User{
		ID:           uuid.NewString(),
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Email:        request.Email,
		PasswordHash: string(hash),
		IsActive:     true,
	}

	return &user, nil
}

func (s *userService) FindByID(ctx context.Context, id string) (*model.User, error) {
	const op = "userService.FindByID"

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.Wrap(apperrors.CategoryNotFound, err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return nil, apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return user, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "userService.FindByEmail"

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.Wrap(apperrors.CategoryNotFound, err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return nil, apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id string, request *model.UpdateUserRequest) error {
	const op = "userService.Update"

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return apperrors.Wrap(apperrors.CategoryNotFound, err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	user.FirstName = *request.FirstName
	user.LastName = *request.LastName

	if err = s.userRepo.Update(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return apperrors.Wrap(apperrors.CategoryNotFound, err.Error(), fmt.Errorf("%s: %w", op, err))
		}
		return apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	const op = "userService.Delete"

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return apperrors.Wrap(apperrors.CategoryInternal, err.Error(), fmt.Errorf("%s: %w", op, err))
	}

	return nil
}
