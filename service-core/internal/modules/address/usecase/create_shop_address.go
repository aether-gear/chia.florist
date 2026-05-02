package usecase

import (
	"fmt"
	"time"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
)

type CreateShopAddressUsecase struct {
	shopAddressRepo repository.ShopAddressRepository
}

func NewCreateShopAddressUsecase(shopAddressRepo repository.ShopAddressRepository) *CreateShopAddressUsecase {
	return &CreateShopAddressUsecase{
		shopAddressRepo: shopAddressRepo,
	}
}

type CreateShopAddressInput struct {
	ShopID      uuid.UUID
	Label       string
	Phone       *string
	IsActive    *bool
	ProvinceID  string
	CityID      string
	DistrictID  string
	VillageID   string
	FullAddress string
	PostalCode  string
}

func (u *CreateShopAddressUsecase) Execute(input CreateShopAddressInput) error {
	var isDefault bool
	if *input.IsActive {
		isDefault = *input.IsActive
	} else {
		isDefault = false
	}

	address := domain.ShopAddress{
		ID:       uuid.New(),
		ShopID:   input.ShopID,
		Label:    input.Label,
		Phone:    input.Phone,
		IsActive: isDefault,
		Detail: domain.AddressDetail{
			ProvinceID:  input.ProvinceID,
			CityID:      input.CityID,
			DistrictID:  input.DistrictID,
			VillageID:   input.VillageID,
			FullAddress: input.FullAddress,
			PostalCode:  input.PostalCode,
		},
		CreatedAt: time.Now(),
	}

	err := u.shopAddressRepo.Create(address)
	if err != nil {
		return fmt.Errorf("failed to save address: %w", err)
	}

	return nil
}
