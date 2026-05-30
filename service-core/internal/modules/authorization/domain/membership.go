package domain

import "github.com/google/uuid"

type MerchantMembership struct {
	MerchantID uuid.UUID
	AccountID  uuid.UUID
	RoleID     uuid.UUID
}
