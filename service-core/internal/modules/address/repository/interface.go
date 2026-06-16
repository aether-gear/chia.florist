package repository

import (
	"context"

	"service-core/internal/modules/address/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type UserAddressRepository interface {
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		addressID uuid.UUID,
	) (*domain.Address, error)

	GetDefaultByUserID(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) (*domain.Address, error)

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

	GetDefaultByShopID(
		ctx context.Context,
		exec transaction.Executor,
		shopID uuid.UUID,
	) (*domain.ShopAddress, error)

	// GetDefaultsByShopIDs retrieves default address grouped by shop IDs.
	// The returned map uses the shop ID as the key and all associated
	// default address record as the value.
	//
	// Example:
	//
	//	shopIDs := []uuid.UUID{shopA, shopB}
	//
	//	result := map[uuid.UUID][]domain.ShopAddress{
	//		shopA: addressDefault,
	//		shopB: addressDefault,
	//	}
	//
	// This allows callers to efficiently look up default addresses
	// belonging to a specific shop without additional filtering.
	GetDefaultsByShopIDs(
		ctx context.Context,
		exec transaction.Executor,
		shopIDs []uuid.UUID,
	) (map[uuid.UUID]domain.ShopAddress, error)

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
