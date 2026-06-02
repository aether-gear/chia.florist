package repository

import (
	"service-core/internal/modules/merchant/domain"

	"github.com/google/uuid"
)

type MerchantRepository interface {
	GetByAccountID(accountID uuid.UUID) (*domain.Merchant, error)
}
