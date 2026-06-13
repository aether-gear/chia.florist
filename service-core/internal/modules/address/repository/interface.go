package repository

import (
	"context"

	"service-core/internal/modules/address/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type UserAddressRepository interface {
	ListByUserID(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) ([]domain.Address, error)

	CountByUserID(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) (*int, error)

	UnsetDefaultByUserID(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) error

	Save(
		ctx context.Context,
		exec transaction.Executor,
		address domain.Address,
	) error

	Delete(
		ctx context.Context,
		exec transaction.Executor,
		addressID uuid.UUID,
	) error
}

type ShopAddressRepository interface {
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.ShopAddress, error)
	// GetDefaultByShopID(shopID string) (*domain.ShopAddress, error)

	FindByShopID(
		ctx context.Context,
		exec transaction.Executor,
		shopID uuid.UUID,
	) ([]domain.ShopAddress, error)

	Create(
		ctx context.Context,
		exec transaction.Executor,
		address domain.ShopAddress,
	) error

	// Update(address domain.ShopAddress) error

	// SetDefault(shopID string, addressID string) error
}
