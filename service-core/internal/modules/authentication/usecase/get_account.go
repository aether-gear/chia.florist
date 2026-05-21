package usecase

import (
	"fmt"

	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/google/uuid"
)

type GetAccountUsecase struct {
	authRepo repository.AccountRepository
}

func NewGetAccountUsecase(
	authRepo repository.AccountRepository,
) *GetAccountUsecase {
	return &GetAccountUsecase{
		authRepo: authRepo,
	}
}

func (u *GetAccountUsecase) Execute(id uuid.UUID) (*domain.Account, error) {
	acc, err := u.authRepo.GetByUserID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	return acc, nil
}
