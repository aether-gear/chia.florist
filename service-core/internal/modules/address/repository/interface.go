package repository

import (
	"context"

	"service-core/internal/modules/address/domain"

	"github.com/google/uuid"
)

type UserAddressRepository interface {
	GetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.Address, error)
	// GetDefault(userID uuid.UUID) (*domain.Address, error)
	Create(
		ctx context.Context,
		address domain.Address,
	) error
}

type ShopAddressRepository interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.ShopAddress, error)
	// GetDefaultByShopID(shopID string) (*domain.ShopAddress, error)

	FindByShopID(
		ctx context.Context,
		shopID uuid.UUID,
	) ([]domain.ShopAddress, error)

	Create(
		ctx context.Context,
		address domain.ShopAddress,
	) error

	// Update(address domain.ShopAddress) error

	// SetDefault(shopID string, addressID string) error
}
