package domain

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	ShopID    uuid.UUID
	Stock     int
	Reserved  int

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (i *Inventory) Validate() error {
	if i.Stock < 0 {
		return ErrInvalidStock
	}

	if i.Reserved < 0 {
		return ErrInvalidReserved
	}

	return nil
}

func (i Inventory) Available() int {
	available := i.Stock - i.Reserved
	if available < 0 {
		return 0
	}

	return available
}
