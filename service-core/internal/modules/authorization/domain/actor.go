package domain

import (
	"service-core/internal/modules/authentication/domain"

	"github.com/google/uuid"
)

type Actor struct {
	AccountID uuid.UUID

	Type domain.AccountType

	MerchantID *uuid.UUID

	Roles       []Role
	Permissions []Permission
}
