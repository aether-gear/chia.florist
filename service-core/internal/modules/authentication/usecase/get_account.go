package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/google/uuid"
)

type GetAccountUsecase struct {
	accountRepo repository.AccountRepository
}

func NewGetAccountUsecase(
	accountRepo repository.AccountRepository,
) *GetAccountUsecase {
	return &GetAccountUsecase{
		accountRepo: accountRepo,
	}
}

func (u *GetAccountUsecase) Execute(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	acc, err := u.accountRepo.GetByUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	return acc, nil
}
