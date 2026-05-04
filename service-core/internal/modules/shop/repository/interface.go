package repository

import (
	"service-core/internal/modules/shop/domain"

	"github.com/google/uuid"
)

type ShopRepository interface {
	GetByID(id uuid.UUID) (*domain.Shop, error)
	// GetByIDs(ids []uuid.UUID) ([]domain.Shop, error)
	// GetActive() ([]domain.Shop, error)

	Create(shop domain.Shop) error
	// Update(shop domain.Shop) error

	// GetSupportedCouriers(shopID uuid.UUID) ([]Courier, error)
}
