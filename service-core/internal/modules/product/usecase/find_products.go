package usecase

import (
	"fmt"

	inventoryR "service-core/internal/modules/inventory/repository"
	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type FindProductsUsecase struct {
	productRepo   repository.ProductRepository
	inventoryRepo inventoryR.InventoryRepository
}

func NewFindProductsUsecase(
	productRepo repository.ProductRepository,
	inventoryRepo inventoryR.InventoryRepository,
) *FindProductsUsecase {
	return &FindProductsUsecase{
		productRepo:   productRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (u *FindProductsUsecase) Execute(
	params repository.FindProductParams,
) (
	[]repository.ProductWithInventory,
	int,
	error,
) {
	products, total, err := u.productRepo.FindProducts(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load products: %w", err)
	}
	if len(products) == 0 {
		return []repository.ProductWithInventory{}, total, nil
	}

	productIDs := make([]uuid.UUID, 0, len(products))
	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}

	inventoryMap, err := u.inventoryRepo.ListByProducts(productIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load inventory for products: %w", err)
	}

	results := make([]repository.ProductWithInventory, 0, len(products))
	for _, product := range products {
		inventories := inventoryMap[product.ID]

		result := repository.ProductWithInventory{
			Product:         product,
			ShopInventories: inventories,
		}

		for _, inventory := range inventories {
			result.Inventory.Stock += inventory.Stock
			result.Inventory.ReservedStock += inventory.Reserved
		}

		results = append(results, result)
	}

	return results, total, nil
}
