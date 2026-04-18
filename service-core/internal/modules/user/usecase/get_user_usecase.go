package usecase

import (
	"service-core/internal/modules/user/domain"
	"service-core/internal/modules/user/repository"

	"github.com/google/uuid"
)

type GetUserUsecase struct {
	userRepo repository.UserRepository
}

func NewGetUserUsecase(userRepo repository.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{
		userRepo: userRepo,
	}
}

func (u *GetUserUsecase) ByID(id uuid.UUID) (*domain.User, error) {
	res, err := u.userRepo.GetByID(id)
	if err != nil {
		return &domain.User{}, err
	}

	return res, nil
}
