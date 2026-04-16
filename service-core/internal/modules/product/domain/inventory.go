package domain

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID            uuid.UUID
	ProductID     uuid.UUID
	Stock         int
	ReservedStock int

	CreatedAt time.Time
	UpdatedAt *time.Time
}
