package usecase

import (
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

func (u *GetAccountUsecase) Execute(id uuid.UUID) (*domain.Account, error) {
	acc, err := u.accountRepo.GetByUserID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	return acc, nil
}
