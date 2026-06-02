package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	AccountStatus string
	AccountType   string
)

const (
	AccountPending   AccountStatus = "pending"
	AccountActive    AccountStatus = "active"
	AccountSuspended AccountStatus = "suspended"
	AccountLocked    AccountStatus = "locked"
)

const (
	AccountTypeCustomer AccountType = "customer"
	AccountTypeMerchant AccountType = "merchant"
)

type Account struct {
	ID     uuid.UUID
	UserID uuid.UUID

	Email    string
	Password string

	Status AccountStatus
	Type   AccountType

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
