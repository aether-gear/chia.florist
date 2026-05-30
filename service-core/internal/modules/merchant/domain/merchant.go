package domain

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID        uuid.UUID
	AccountID uuid.UUID

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
