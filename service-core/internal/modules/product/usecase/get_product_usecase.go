package usecase

import (
	"fmt"

	inventoryR "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetProductUsecase struct {
	repo          repository.ProductRepository
	inventoryRepo inventoryR.InventoryRepository
}

func NewGetProductsUsecase(
	repo repository.ProductRepository,
	inventoryRepo inventoryR.InventoryRepository,
) *GetProductUsecase {
	return &GetProductUsecase{
		repo:          repo,
		inventoryRepo: inventoryRepo,
	}
}

func (u *GetProductUsecase) Execute(id uuid.UUID) (*repository.ProductWithInventory, error) {
	product, err := u.repo.GetByID(id)
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
		result.Inventory.Stock += inventory.Stock
		result.Inventory.ReservedStock += inventory.Reserved
	}

	return &result, nil
}
