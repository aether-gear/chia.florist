package repository

import (
	"context"

	"service-core/internal/modules/shop/domain"

	"github.com/google/uuid"
)

type ShopRepository interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Shop, error)
	// GetByIDs(ids []uuid.UUID) ([]domain.Shop, error)
	// GetActive() ([]domain.Shop, error)

	Create(
		ctx context.Context,
		shop domain.Shop,
	) error
	// Update(shop domain.Shop) error

	// GetSupportedCouriers(shopID uuid.UUID) ([]Courier, error)
}
