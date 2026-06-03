package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
)

type GetShopAddressUsecase struct {
	shopAddressRepo repository.ShopAddressRepository
}

func NewGetShopAddressUsecase(
	shopAddressRepo repository.ShopAddressRepository,
) *GetShopAddressUsecase {
	return &GetShopAddressUsecase{
		shopAddressRepo: shopAddressRepo,
	}
}

func (u *GetShopAddressUsecase) GetByID(ctx context.Context, addressID uuid.UUID) (*domain.ShopAddress, error) {
	res, err := u.shopAddressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
