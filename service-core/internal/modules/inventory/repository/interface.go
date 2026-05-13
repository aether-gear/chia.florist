package repository

import (
	"service-core/internal/modules/inventory/domain"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	GetByProductIDAndShopID(productID uuid.UUID, shopID uuid.UUID) (*domain.Inventory, error)

	ListByProductID(productID uuid.UUID) ([]domain.Inventory, error)
	ListByProductIDs(productIDs []uuid.UUID) (map[uuid.UUID][]domain.Inventory, error)
	ListByShopID(shopID uuid.UUID) ([]domain.Inventory, error)

	Create(inventory *domain.Inventory) error
}
