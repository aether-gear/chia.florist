package domain

import "github.com/google/uuid"

type RoleCode string

const (
	// RoleStaffAdmin grants full access to staff resources and
	// administrative operations within a staff workspace.
	RoleStaffAdmin RoleCode = "admin"

	// RoleStaff grants operational access to staff resources
	// with permissions restricted by business policy.
	RoleStaff RoleCode = "staff"
)

// Role defines a collection of permissions that can be assigned to an
// account through a staff membership.
//
// Roles provide a reusable mechanism for grouping permissions and
// implementing role-based access control (RBAC) within a staff
// workspace.
type Role struct {
	ID uuid.UUID

	Code RoleCode
	Name string
}
