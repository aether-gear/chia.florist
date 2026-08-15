package domain

import (
	"time"

	"github.com/google/uuid"
)

// StaffMembership represents the assignment of an account to a
// staff under a specific role.
//
// A membership establishes the relationship between an account and a
// staff and serves as the source of role-based access control within
// a staff workspace.
type StaffMembership struct {
	ID uuid.UUID

	StaffID   uuid.UUID
	AccountID uuid.UUID
	RoleID    uuid.UUID

	CreatedBy uuid.UUID

	CreatedAt time.Time
}

// StaffAccountMember represents account and user details bound to a staff unit.
type StaffAccountMember struct {
	AccountID   uuid.UUID
	UserID      uuid.UUID
	Email       string
	Name        string
	Username    string
	Phone       *string
	AvatarURL   *string
	Role        Role
	LastLoginAt *time.Time
	CreatedAt   time.Time
}

