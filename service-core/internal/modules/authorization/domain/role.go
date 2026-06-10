package domain

import "github.com/google/uuid"

const (
	RoleMerchantAdmin = "merchant_admin"
	RoleMerchantStaff = "merchant_staff"
)

type Role struct {
	ID uuid.UUID

	Code string
	Name string
}
