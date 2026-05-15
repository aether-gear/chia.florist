package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountPendingVerification AccountStatus = "pending_verification"
	AccountActive              AccountStatus = "active"
	AccountSuspended           AccountStatus = "suspended"
	AccountLocked              AccountStatus = "locked"
)

type Account struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Email    string
	Password string

	Status AccountStatus

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
