package usecase

import (
	"fmt"

	inventoryR "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetProductUsecase struct {
	productRepo   repository.ProductRepository
	inventoryRepo inventoryR.InventoryRepository
}

func NewGetProductUsecase(
	productRepo repository.ProductRepository,
	inventoryRepo inventoryR.InventoryRepository,
) *GetProductUsecase {
	return &GetProductUsecase{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (u *GetProductUsecase) Execute(id uuid.UUID) (*repository.ProductWithInventory, error) {
	product, err := u.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load products with inventory: %w", err)
	}
	if product == nil {
		return nil, nil
	}

	inventories, err := u.inventoryRepo.ListByProduct(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load inventory by product: %w", err)
	}

	result := repository.ProductWithInventory{
		Product:         *product,
		ShopInventories: inventories,
	}

	for _, inventory := range inventories {
		result.Inventory.TotalStock += inventory.TotalStock
		result.Inventory.ReservedStock += inventory.ReservedStock
	}

	return &result, nil
}
