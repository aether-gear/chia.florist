package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentAccount struct {
	ID       uuid.UUID
	MethodID uuid.UUID

	AccountName   string
	AccountNumber *string
	PhoneNumber   string
	QRString      *string

	IsActive bool

	CurrentLoad int
	LastUsedAt  *time.Time

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
