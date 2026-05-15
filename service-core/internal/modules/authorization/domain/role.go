package domain

import "github.com/google/uuid"

const (
	RoleAdmin    = "admin"
	RoleSeller   = "seller"
	RoleCustomer = "customer"
)

type Role struct {
	ID uuid.UUID

	Code string
	Name string
}
