package repository

import "service-core/internal/modules/shop/domain"

type ShopRepository interface {
	// GetByID(id string) (*domain.Shop, error)
	// GetActive() ([]domain.Shop, error)

	Create(shop domain.Shop) error
	// Update(shop domain.Shop) error
}
