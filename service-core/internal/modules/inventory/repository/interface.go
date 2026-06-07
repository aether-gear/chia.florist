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
