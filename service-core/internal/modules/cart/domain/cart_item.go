package domain

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	ShopID    uuid.UUID
	Quantity  int
	DeletedAt *time.Time
}
