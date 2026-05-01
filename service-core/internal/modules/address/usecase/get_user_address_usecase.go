package usecase

import (
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
)

type GetAddressUsecase struct {
	addressRepo repository.UserAddressRepository
}

func NewGetAddressUsecase(
	addressRepo repository.UserAddressRepository,
) *GetAddressUsecase {
	return &GetAddressUsecase{
		addressRepo: addressRepo,
	}
}

func (u *GetAddressUsecase) GetByUserID(userID uuid.UUID) ([]domain.Address, error) {
	res, err := u.addressRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
