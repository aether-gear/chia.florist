package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ItemType discriminates standard catalog products from
// custom-designed flower boards.
type ItemType string

const (
	ItemTypeStandard ItemType = "standard"
	ItemTypeCustom   ItemType = "custom"
)

// CartItem represents one line in a customer's cart.
//
// For standard items, ProductID is non-nil and CustomDesign is nil.
// For custom items, ProductID is nil and CustomDesign holds the
type CartItem struct {
	ID           uuid.UUID
	ItemType     ItemType
	ProductID    *uuid.UUID
	ShopID       uuid.UUID
	Quantity     int
	CustomDesign json.RawMessage
	DeletedAt    *time.Time
}
