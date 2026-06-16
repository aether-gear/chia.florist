package domain

import "github.com/google/uuid"

type RoleCode string

const (
	// RoleMerchantAdmin grants full access to merchant resources and
	// administrative operations within a merchant workspace.
	RoleMerchantAdmin RoleCode = "merchant_admin"

	// RoleMerchantStaff grants operational access to merchant resources
	// with permissions restricted by business policy.
	RoleMerchantStaff RoleCode = "merchant_staff"
)

// Role defines a collection of permissions that can be assigned to an
// account through a merchant membership.
//
// Roles provide a reusable mechanism for grouping permissions and
// implementing role-based access control (RBAC) within a merchant
// workspace.
type Role struct {
	ID uuid.UUID

	Code RoleCode
	Name string
}
