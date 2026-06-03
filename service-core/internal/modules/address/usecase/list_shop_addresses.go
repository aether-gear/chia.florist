package usecase

import (
	"context"
	"fmt"

	"service-core/internal/modules/address/domain"
	"service-core/internal/modules/address/repository"

	"github.com/google/uuid"
)

type ListShopAddressesUsecase struct {
	shopAddressRepo repository.ShopAddressRepository
}

func NewListShopAddressesUsecase(
	shopAddressRepo repository.ShopAddressRepository,
) *ListShopAddressesUsecase {
	return &ListShopAddressesUsecase{
		shopAddressRepo: shopAddressRepo,
	}
}

func (u *ListShopAddressesUsecase) FindByShopID(
	ctx context.Context,
	shopID uuid.UUID,
) ([]domain.ShopAddress, error) {
	res, err := u.shopAddressRepo.FindByShopID(ctx, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve address: %w", err)
	}

	return res, nil
}
