package domain

import (
	"time"

	"github.com/google/uuid"
)

type MerchantMembership struct {
	ID uuid.UUID

	MerchantID uuid.UUID
	AccountID  uuid.UUID
	RoleID     uuid.UUID

	CreatedBy uuid.UUID

	CreatedAt time.Time
}
