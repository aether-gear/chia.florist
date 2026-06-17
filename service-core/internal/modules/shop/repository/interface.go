package repository

import (
	"context"

	"service-core/internal/modules/shop/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ShopRepository interface {
	FindByParams(
		ctx context.Context,
		exec transaction.Executor,
		params FindShopsParams,
	) ([]domain.Shop, int, error)

	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Shop, error)
	// GetByIDs(ids []uuid.UUID) ([]domain.Shop, error)
	// GetActive() ([]domain.Shop, error)

	Create(
		ctx context.Context,
		exec transaction.Executor,
		shop domain.Shop,
	) error
	// Update(shop domain.Shop) error

	// GetSupportedCouriers(shopID uuid.UUID) ([]Courier, error)
}
