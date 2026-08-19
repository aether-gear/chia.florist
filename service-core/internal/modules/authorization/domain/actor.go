package domain

import (
	"service-core/internal/modules/authentication/domain"

	"github.com/google/uuid"
)

// Actor represents the authorization identity of an account within the
// application.
//
// An Actor is derived from an authenticated account and contains the
// roles and permissions granted to that account within a staff scope.
// It is primarily used by authorization and access control workflows to
// determine what actions an authenticated account is allowed to perform.
//
// An Actor does not represent a persisted entity and is typically
// constructed from account, membership, role, and permission data.
type Actor struct {
	AccountID uuid.UUID

	Type domain.AccountType

	StaffID *uuid.UUID

	Roles []Role

	Permissions map[uuid.UUID][]string
	Rules       map[uuid.UUID]map[string]any
}

func (a *Actor) HasRole(role RoleCode) bool {
	for _, r := range a.Roles {
		if r.Code == role {
			return true
		}
	}

	return false
}

func (a *Actor) IsSuperAdmin() bool {
	return a.HasRole(RoleStaffAdmin)
}

func (a *Actor) HasPermission(shopID uuid.UUID, permission string) bool {
	if a.IsSuperAdmin() {
		return true
	}
	if a.Permissions == nil {
		return false
	}
	perms, exists := a.Permissions[shopID]
	if !exists {
		return false
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

func (a *Actor) GetAssignedShopIDs() []uuid.UUID {
	if a.Permissions == nil {
		return nil
	}
	ids := make([]uuid.UUID, 0, len(a.Permissions))
	for shopID := range a.Permissions {
		ids = append(ids, shopID)
	}
	return ids
}
