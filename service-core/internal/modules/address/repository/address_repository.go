package repository

import (
	"service-core/internal/modules/address/domain"

	"github.com/google/uuid"
)

type AddressRepository interface {
	GetByUserID(userID uuid.UUID) ([]domain.Address, error)
	Save(address domain.Address) error
}
