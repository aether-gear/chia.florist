package repository

import (
	"service-core/internal/modules/product/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByID(id uuid.UUID) (*ProductWithInventory, error)

	FindProducts(params FindProductParams) ([]ProductWithInventory, int, error)
	FindByIDs(IDs []uuid.UUID) ([]ProductWithInventory, error)

	CreateProduct(product *domain.Product) error
	CreateInventory(inventory *domain.Inventory) error
}

// type ProductAvailabilityRepository interface {
// 	GetAvailableShops(productID uuid.UUID) ([]domain.ProductShopAvailability, error)
// 	GetProductsByShop(shopID uuid.UUID, productIDs []uuid.UUID) ([]domain.ProductShopAvailability, error)
// }

type ProductInventoryRepository interface {
	CreateInventory(inventory *domain.Inventory) error
}
