package usecase

import (
	"fmt"

	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"

	"github.com/google/uuid"
)

type GetShopUsecase struct {
	shopRepo repository.ShopRepository
}

func NewGetShopUsecase(shopRepo repository.ShopRepository) *GetShopUsecase {
	return &GetShopUsecase{
		shopRepo: shopRepo,
	}
}

func (u *GetShopUsecase) GetByID(shopID uuid.UUID) (*domain.Shop, error) {
	shop, err := u.shopRepo.GetByID(shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shop: %w", err)
	}

	return shop, nil
}
