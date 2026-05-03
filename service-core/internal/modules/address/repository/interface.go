package repository

import (
	"service-core/internal/modules/address/domain"

	"github.com/google/uuid"
)

type UserAddressRepository interface {
	GetByUserID(userID uuid.UUID) ([]domain.Address, error)
	Create(address domain.Address) error
}

type ShopAddressRepository interface {
	GetByID(id uuid.UUID) (*domain.ShopAddress, error)
	FindByShopID(shopID uuid.UUID) ([]domain.ShopAddress, error)

	// GetDefaultByShopID(shopID string) (*domain.ShopAddress, error)

	Create(address domain.ShopAddress) error
	// Update(address domain.ShopAddress) error

	// SetDefault(shopID string, addressID string) error
}
