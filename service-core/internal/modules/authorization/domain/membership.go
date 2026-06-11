package domain

import (
	"time"

	"github.com/google/uuid"
)

// MerchantMembership represents the assignment of an account to a
// merchant under a specific role.
//
// A membership establishes the relationship between an account and a
// merchant and serves as the source of role-based access control within
// a merchant workspace.
type MerchantMembership struct {
	ID uuid.UUID

	MerchantID uuid.UUID
	AccountID  uuid.UUID
	RoleID     uuid.UUID

	CreatedBy uuid.UUID

	CreatedAt time.Time
}
