package usecase

import (
	"service-core/internal/features/auth/domain"
	"service-core/internal/features/auth/repository"

	"github.com/google/uuid"
)

type GetAccountUsecase struct {
	repo repository.AuthRepository
}

func NewGetAccountUsecase(repo repository.AuthRepository) *GetAccountUsecase {
	return &GetAccountUsecase{repo: repo}
}

func (u *GetAccountUsecase) ById(id uuid.UUID) (*domain.Account, error) {
	return u.repo.GetByID(id)
}
