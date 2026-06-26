package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents shared profile information used across account types.
// A user may be associated with customer or staff roles through
// an authenticated account.
type User struct {
	ID       uuid.UUID
	Name     string
	Username string
	Phone    *string

	AvatarURL *string

	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	LastLoginAt *time.Time
}
