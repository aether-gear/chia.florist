package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/user/domain"
	"service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetUserUsecase struct {
	userRepo repository.UserRepository
	executor transaction.Executor
}

func NewGetUserUsecase(
	userRepo repository.UserRepository,
	executor transaction.Executor,
) *GetUserUsecase {
	return &GetUserUsecase{
		userRepo: userRepo,
		executor: executor,
	}
}

func (u *GetUserUsecase) ByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {
	user, err := u.userRepo.GetByID(ctx, u.executor, id)
	if err != nil {
		return &domain.User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return user, nil
}
