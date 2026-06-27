package domain

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID uuid.UUID

	CustomerID uuid.UUID

	Items []CartItem

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (c *Cart) AddItem(productID uuid.UUID, shopID uuid.UUID, qty int) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	for i := range c.Items {
		if c.Items[i].ProductID == productID && c.Items[i].ShopID == shopID {
			c.Items[i].Quantity += qty
			return nil
		}
	}

	c.Items = append(c.Items, CartItem{
		ID:        uuid.New(),
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  qty,
	})

	return nil
}

func (c *Cart) SetItem(productID uuid.UUID, shopID uuid.UUID, qty int) error {
	if qty < 0 {
		return ErrInvalidQuantity
	}

	for i := range c.Items {
		if c.Items[i].ProductID == productID && c.Items[i].ShopID == shopID {

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
		ShopID:    shopID,
		Quantity:  qty,
		DeletedAt: nil,
	})

	return nil
}

func (c *Cart) RemoveItem(productID uuid.UUID, shopID uuid.UUID) {
	for i := range c.Items {
		if c.Items[i].ProductID == productID && c.Items[i].ShopID == shopID {
			now := time.Now()
			c.Items[i].DeletedAt = &now
			return
		}
	}
}

func (c *Cart) HasItem(productID uuid.UUID, shopID uuid.UUID) bool {
	for i := range c.Items {
		if c.Items[i].ProductID == productID && c.Items[i].ShopID == shopID {
			return true
		}
	}
	return false
}

func (c *Cart) FindItem(productID uuid.UUID, shopID uuid.UUID) *CartItem {
	for i := range c.Items {
		if c.Items[i].ProductID == productID && c.Items[i].ShopID == shopID {
			return &c.Items[i]
		}
	}

	return nil
}

func (c *Cart) HasProductInAnotherShop(productID uuid.UUID, shopID uuid.UUID) bool {
	for _, item := range c.Items {
		if item.DeletedAt != nil {
			continue
		}

		if item.ProductID == productID && item.ShopID != shopID {
			return true
		}
	}

	return false
}
