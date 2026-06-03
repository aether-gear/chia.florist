package repository

import (
	"context"

	"service-core/internal/modules/inventory/domain"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	GetByProductIDAndShopID(
		ctx context.Context,
		productID uuid.UUID,
		shopID uuid.UUID,
	) (*domain.Inventory, error)

	ListByProductID(
		ctx context.Context,
		productID uuid.UUID,
	) ([]domain.Inventory, error)
	ListByProductIDs(
		ctx context.Context,
		productIDs []uuid.UUID,
	) (map[uuid.UUID][]domain.Inventory, error)
	ListByShopID(
		ctx context.Context,
		shopID uuid.UUID,
	) ([]domain.Inventory, error)

	Create(
		ctx context.Context,
		inventory *domain.Inventory,
	) error
}
