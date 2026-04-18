package usecase

import (
	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
)

type GetAddressUsecase struct {
	addressRepo repository.AddressRepository
}

func NewGetAddressUsecase(addressRepo repository.AddressRepository) *GetAddressUsecase {
	return &GetAddressUsecase{
		addressRepo: addressRepo,
	}
}

func (u *GetAddressUsecase) GetByUserID(userID uuid.UUID) ([]domain.Address, error) {
	res, err := u.addressRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
