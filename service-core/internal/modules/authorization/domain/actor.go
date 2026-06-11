package domain

import (
	"service-core/internal/modules/authentication/domain"

	"github.com/google/uuid"
)

// Actor represents the authorization identity of an account within the
// application.
//
// An Actor is derived from an authenticated account and contains the
// roles and permissions granted to that account within a merchant scope.
// It is primarily used by authorization and access control workflows to
// determine what actions an authenticated account is allowed to perform.
//
// An Actor does not represent a persisted entity and is typically
// constructed from account, membership, role, and permission data.
type Actor struct {
	AccountID uuid.UUID

	Type domain.AccountType

	MerchantID *uuid.UUID

	Roles       []Role
	Permissions []Permission
}
