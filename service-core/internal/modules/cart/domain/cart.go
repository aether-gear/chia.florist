package domain

import (
	"encoding/json"
	"time"

	appclock "service-core/internal/common/clock"

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
		item := &c.Items[i]
		// Skip custom items,
		// they are never deduplicated by product_id
		if item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
			continue
		}
		if *item.ProductID == productID && item.ShopID == shopID {
			item.Quantity += qty
			return nil
		}
	}

	pid := productID
	c.Items = append(c.Items, CartItem{
		ID:                 uuid.New(),
		ProductVariantType: ProductVariantTypeStandard,
		ProductID:          &pid,
		ShopID:             shopID,
		Quantity:           qty,
	})

	return nil
}

// AddCustomItem appends a new custom design
// item to the cart.
//
// Custom items are never deduplicated,
// each canvas design is independent.
func (c *Cart) AddCustomItem(shopID uuid.UUID, qty int, design json.RawMessage) error {
	if qty <= 0 {
		return ErrInvalidQuantity
	}
	c.Items = append(c.Items, CartItem{
		ID:                 uuid.New(),
		ProductVariantType: ProductVariantTypeCustom,
		ProductID:          nil,
		ShopID:             shopID,
		Quantity:           qty,
		CustomDesign:       design,
	})
	return nil
}

func (c *Cart) SetItem(productID uuid.UUID, shopID uuid.UUID, qty int) error {
	if qty < 0 {
		return ErrInvalidQuantity
	}

	for i := range c.Items {
		item := &c.Items[i]
		// Skip custom items,
		// they are managed by ID, not by product_id
		if item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
			continue
		}
		if *item.ProductID == productID && item.ShopID == shopID {
			if qty == 0 {
				now := appclock.Now()
				item.DeletedAt = &now
				return nil
			}

			item.Quantity = qty
			item.DeletedAt = nil
			return nil
		}
	}

	if qty == 0 {
		return nil
	}

	pid := productID
	c.Items = append(c.Items, CartItem{
		ID:                 uuid.New(),
		ProductVariantType: ProductVariantTypeStandard,
		ProductID:          &pid,
		ShopID:             shopID,
		Quantity:           qty,
		DeletedAt:          nil,
	})

	return nil
}

func (c *Cart) RemoveItem(productID uuid.UUID, shopID uuid.UUID) bool {
	for i := range c.Items {
		item := &c.Items[i]
		if item.DeletedAt != nil || item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
			continue
		}
		if *item.ProductID == productID && item.ShopID == shopID {
			now := appclock.Now()
			item.DeletedAt = &now
			return true
		}
	}
	return false
}

// RemoveItemByID soft-deletes any cart item (standard or custom) by its own UUID.
// Returns false if the item was not found or already deleted.
func (c *Cart) RemoveItemByID(cartItemID uuid.UUID) bool {
	for i := range c.Items {
		if c.Items[i].ID == cartItemID && c.Items[i].DeletedAt == nil {
			now := appclock.Now()
			c.Items[i].DeletedAt = &now
			return true
		}
	}
	return false
}

// RemoveProduct soft-deletes an active standard cart item matching productID regardless of shopID.
// Returns false if the item was not found or already deleted.
func (c *Cart) RemoveProduct(productID uuid.UUID) bool {
	for i := range c.Items {
		item := &c.Items[i]
		if item.DeletedAt != nil || item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
			continue
		}
		if *item.ProductID == productID {
			now := appclock.Now()
			item.DeletedAt = &now
			return true
		}
	}
	return false
}

// RemoveCustomItem soft-deletes a custom cart item by its own UUID.
// Returns false if the item was not found.
func (c *Cart) RemoveCustomItem(cartItemID uuid.UUID) bool {
	for i := range c.Items {
		if c.Items[i].ID == cartItemID && c.Items[i].ProductVariantType == ProductVariantTypeCustom && c.Items[i].DeletedAt == nil {
			now := appclock.Now()
			c.Items[i].DeletedAt = &now
			return true
		}
	}
	return false
}

func (c *Cart) HasItem(productID uuid.UUID, shopID uuid.UUID) bool {
	for i := range c.Items {
		item := &c.Items[i]
		if item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
			continue
		}
		if *item.ProductID == productID && item.ShopID == shopID {
			return true
		}
	}
	return false
}

func (c *Cart) FindItem(productID uuid.UUID, shopID uuid.UUID) *CartItem {
	for i := range c.Items {
		item := &c.Items[i]
		if item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
			continue
		}
		if *item.ProductID == productID && item.ShopID == shopID {
			return item
		}
	}

	return nil
}

func (c *Cart) HasProductInAnotherShop(productID uuid.UUID, shopID uuid.UUID) bool {
	for _, item := range c.Items {
		if item.DeletedAt != nil {
			continue
		}
		// Custom items have no product_id to compare
		if item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
			continue
		}
		if *item.ProductID == productID && item.ShopID != shopID {
			return true
		}
	}

	return false
}

// ChangeItemShop updates the fulfillment shop for a cart item identified by cartItemID.
// For standard products, if an item with the same product ID already exists in newShopID,
// their quantities are merged and the original item is soft-deleted.
func (c *Cart) ChangeItemShop(cartItemID uuid.UUID, newShopID uuid.UUID) bool {
	var targetIdx = -1
	for i := range c.Items {
		if c.Items[i].ID == cartItemID && c.Items[i].DeletedAt == nil {
			targetIdx = i
			break
		}
	}

	if targetIdx == -1 {
		return false
	}

	item := &c.Items[targetIdx]
	if item.ShopID == newShopID {
		return true
	}

	// Custom items have no product ID deduplication; directly update shop_id
	if item.ProductVariantType == ProductVariantTypeCustom || item.ProductID == nil {
		item.ShopID = newShopID
		return true
	}

	// Standard product item: check if product already exists in newShopID
	productID := *item.ProductID
	for i := range c.Items {
		if i != targetIdx && c.Items[i].DeletedAt == nil &&
			c.Items[i].ProductVariantType != ProductVariantTypeCustom &&
			c.Items[i].ProductID != nil &&
			*c.Items[i].ProductID == productID &&
			c.Items[i].ShopID == newShopID {

			// Merge quantity into existing item for newShopID and soft-delete target item
			c.Items[i].Quantity += item.Quantity
			now := appclock.Now()
			item.DeletedAt = &now
			return true
		}
	}

	item.ShopID = newShopID
	return true
}
