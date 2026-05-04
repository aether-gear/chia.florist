package repository

import (
	"service-core/internal/modules/inventory/domain"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	GetByProductAndShop(productID uuid.UUID, shopID uuid.UUID) (*domain.Inventory, error)

	ListByProduct(productID uuid.UUID) ([]domain.Inventory, error)
	ListByProducts(productIDs []uuid.UUID) (map[uuid.UUID][]domain.Inventory, error)
	ListByShop(shopID uuid.UUID) ([]domain.Inventory, error)

	Create(inventory *domain.Inventory) error
}
