package domain

import (
	"fmt"
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

func (i *Inventory) Validate() error {
	if i.Stock < 0 {
		return fmt.Errorf("invalid inventory: stock cannot be negative")
	}

	return nil
}
