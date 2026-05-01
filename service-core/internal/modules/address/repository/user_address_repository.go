package repository

import (
	"service-core/internal/modules/address/domain"

	"github.com/google/uuid"
)

type UserAddressRepository interface {
	GetByUserID(userID uuid.UUID) ([]domain.Address, error)
	Create(address domain.Address) error
}
