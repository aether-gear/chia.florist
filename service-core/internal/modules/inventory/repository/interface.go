package repository

import (
	"context"

	"service-core/internal/modules/inventory/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	GetByProductIDAndShopID(
		ctx context.Context,
		exec transaction.Executor,
		productID uuid.UUID,
		shopID uuid.UUID,
	) (*domain.Inventory, error)

	ListByProductID(
		ctx context.Context,
		exec transaction.Executor,
		productID uuid.UUID,
	) ([]domain.Inventory, error)

	// ListByProductIDs retrieves inventories grouped by product IDs.
	// The returned map uses the product ID as the key and all associated
	// inventory records as the value.
	//
	// Example:
	//
	//	productIDs := []uuid.UUID{productA, productB}
	//
	//	result := map[uuid.UUID][]domain.Inventory{
	//		productA: {
	//			inventoryA1,
	//			inventoryA2,
	//		},
	//		productB: {
	//			inventoryB1,
	//		},
	//	}
	//
	// This allows callers to efficiently look up inventories belonging to a
	// specific product without additional filtering.
	ListByProductIDs(
		ctx context.Context,
		exec transaction.Executor,
		productIDs []uuid.UUID,
	) (map[uuid.UUID][]domain.Inventory, error)

	ListByShopID(
		ctx context.Context,
		exec transaction.Executor,
		shopID uuid.UUID,
	) ([]domain.Inventory, error)

	Create(
		ctx context.Context,
		exec transaction.Executor,
		inventory *domain.Inventory,
	) error
}
