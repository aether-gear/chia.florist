package usecase

import (
	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"
	"time"

	"github.com/google/uuid"
)

type CreateAddressInput struct {
	UserID       uuid.UUID
	ReceiverName string
	Phone        *string
	IsDefault    *bool
	ProvinceID   string
	CityID       string
	DistrictID   string
	VillageID    string
	FullAddress  string
	PostalCode   string
}

type CreateAddressUsecase struct {
	addressRepo repository.AddressRepository
}

func NewCreateAddressUsecase(addressRepo repository.AddressRepository) *CreateAddressUsecase {
	return &CreateAddressUsecase{
		addressRepo: addressRepo,
	}
}

func (u *CreateAddressUsecase) Execute(input CreateAddressInput) error {
	var isDefault bool
	if *input.IsDefault {
		isDefault = *input.IsDefault
	} else {
		isDefault = false
	}

	address := domain.Address{
		ID:           uuid.New(),
		UserID:       input.UserID,
		ReceiverName: input.ReceiverName,
		Phone:        input.Phone,
		IsDefault:    isDefault,
		ProvinceID:   input.ProvinceID,
		CityID:       input.CityID,
		DistrictID:   input.DistrictID,
		VillageID:    input.VillageID,
		FullAddress:  input.FullAddress,
		PostalCode:   input.PostalCode,
		CreatedAt:    time.Now(),
	}

	err := u.addressRepo.Save(address)
	if err != nil {
		return err
	}

	return nil
}
