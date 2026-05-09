package domain

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID            uuid.UUID
	ProductID     uuid.UUID
	ShopID        uuid.UUID
	TotalStock    int
	ReservedStock int

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (i *Inventory) Validate() error {
	if i.TotalStock < 0 {
		return ErrInvalidStock
	}

	if i.ReservedStock < 0 {
		return ErrInvalidReserved
	}

	if i.ReservedStock > i.TotalStock {
		return ErrReservedExceedsStock
	}

	return nil
}

func (i Inventory) Available() int {
	return i.TotalStock - i.ReservedStock
}

func (i *Inventory) Reserve(qty int) error {
	if qty <= 0 {
		return ErrInvalidStock
	}

	if i.Available() < qty {
		return ErrInsufficientStock
	}

	i.ReservedStock += qty

	return i.Validate()
}

func (i *Inventory) Release(qty int) error {
	if qty <= 0 {
		return ErrInvalidReserved
	}

	if i.ReservedStock < qty {
		return ErrInsufficientReserved
	}

	i.ReservedStock -= qty

	return i.Validate()
}

func (i *Inventory) Commit(qty int) error {
	if qty <= 0 {
		return ErrInvalidReserved
	}

	if i.ReservedStock < qty {
		return ErrInsufficientReserved
	}

	i.TotalStock -= qty
	i.ReservedStock -= qty

	return i.Validate()
}
