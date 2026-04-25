package usecase

import (
	"fmt"

	"service-core/internal/modules/auth/domain"
	"service-core/internal/modules/auth/repository"

	"github.com/google/uuid"
)

type GetAccountUsecase struct {
	repo repository.AuthRepository
}

func NewGetAccountUsecase(
	repo repository.AuthRepository,
) *GetAccountUsecase {
	return &GetAccountUsecase{repo: repo}
}

func (u *GetAccountUsecase) ById(id uuid.UUID) (*domain.Account, error) {
	acc, err := u.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	return acc, nil
}
