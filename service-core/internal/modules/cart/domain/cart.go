package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID uuid.UUID

	UserID uuid.UUID

	Items []CartItem

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (c *Cart) AddItem(productID uuid.UUID, qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			c.Items[i].Quantity += qty
			return nil
		}
	}

	c.Items = append(c.Items, CartItem{
		ID:        uuid.New(),
		ProductID: productID,
		Quantity:  qty,
	})

	return nil
}

func (c *Cart) SetItem(productID uuid.UUID, qty int) error {
	if qty < 0 {
		return ErrInvalidQuantity
	}

	for i := range c.Items {
		if c.Items[i].ProductID == productID {

			if qty == 0 {
				now := time.Now()
				c.Items[i].DeletedAt = &now
				return nil
			}

			c.Items[i].Quantity = qty
			c.Items[i].DeletedAt = nil
			return nil
		}
	}

	if qty == 0 {
		return nil
	}

	c.Items = append(c.Items, CartItem{
		ID:        uuid.New(),
		ProductID: productID,
		Quantity:  qty,
		DeletedAt: nil,
	})

	return nil
}

func (c *Cart) RemoveItem(productID uuid.UUID) {
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			now := time.Now()
			c.Items[i].DeletedAt = &now
			return
		}
	}
}

func (c *Cart) HasItem(productID uuid.UUID) bool {
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			return true
		}
	}
	return false
}
