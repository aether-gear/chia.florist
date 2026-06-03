package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/user/domain"
	"service-core/internal/modules/user/repository"

	"github.com/google/uuid"
)

type GetUserUsecase struct {
	userRepo repository.UserRepository
}

func NewGetUserUsecase(
	userRepo repository.UserRepository,
) *GetUserUsecase {
	return &GetUserUsecase{
		userRepo: userRepo,
	}
}

func (u *GetUserUsecase) ByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return &domain.User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return user, nil
}
