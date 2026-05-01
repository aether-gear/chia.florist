package usecase

import (
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
)

type FindShopAddressUsecase struct {
	shopAddressRepo repository.ShopAddressRepository
}

func NewFindShopAddressUsecase(
	shopAddressRepo repository.ShopAddressRepository,
) *FindShopAddressUsecase {
	return &FindShopAddressUsecase{
		shopAddressRepo: shopAddressRepo,
	}
}

func (u *FindShopAddressUsecase) FindByShopID(shopID uuid.UUID) ([]domain.ShopAddress, error) {
	res, err := u.shopAddressRepo.FindByShopID(shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
