package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID uuid.UUID

	UserID uuid.UUID

	ReceiverName string
	Phone        *string

	IsDefault bool

	Detail AddressDetail

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
