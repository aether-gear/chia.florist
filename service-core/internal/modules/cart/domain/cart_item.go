package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProductVariantType discriminates standard catalog products from
// custom-designed flower boards.
type ProductVariantType string

const (
	ProductVariantTypeStandard ProductVariantType = "standard"
	ProductVariantTypeCustom   ProductVariantType = "custom"
)

// CartItem represents one line in a customer's cart.
//
// For standard items, ProductID is non-nil and CustomDesign is nil.
// For custom items, ProductID is nil and CustomDesign holds the
type CartItem struct {
	ID                 uuid.UUID
	ProductVariantType ProductVariantType
	ProductID          *uuid.UUID
	ShopID             uuid.UUID
	Quantity           int
	CustomDesign       json.RawMessage
	DeletedAt          *time.Time
}
