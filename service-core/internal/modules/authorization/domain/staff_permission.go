package domain

import (
	"time"

	"github.com/google/uuid"
)

type StaffPermission struct {
	ID      uuid.UUID
	StaffID uuid.UUID

	ShopID   uuid.UUID
	ShopName string

	Permissions []string
	Rules       map[string]any

	CreatedAt time.Time
	UpdatedAt *time.Time
}
