package usecase

import (
	"fmt"

	"service-core/internal/modules/auth/domain"
	"service-core/internal/modules/auth/repository"

	"github.com/google/uuid"
)

type GetAccountUsecase struct {
	authRepo repository.AuthRepository
}

func NewGetAccountUsecase(
	authRepo repository.AuthRepository,
) *GetAccountUsecase {
	return &GetAccountUsecase{
		authRepo: authRepo,
	}
}

func (u *GetAccountUsecase) Execute(id uuid.UUID) (*domain.Account, error) {
	acc, err := u.authRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	return acc, nil
}
