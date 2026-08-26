package domain

import (
	"encoding/json"
	"strings"
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

// ItemOptions captures variant specifications for standard products (size and floral jambul).
type ItemOptions struct {
	Size   string
	Jambul string
}

func (o ItemOptions) Normalized() ItemOptions {
	size := strings.ToLower(strings.TrimSpace(o.Size))
	if size == "" {
		size = "small"
	}
	jambul := strings.ToLower(strings.TrimSpace(o.Jambul))
	if jambul == "" {
		jambul = "none"
	}
	return ItemOptions{
		Size:   size,
		Jambul: jambul,
	}
}

func (o ItemOptions) Equals(other ItemOptions) bool {
	n1 := o.Normalized()
	n2 := other.Normalized()
	return n1.Size == n2.Size && n1.Jambul == n2.Jambul
}

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
	ItemOptions        ItemOptions
	CustomDesign       json.RawMessage
	DeletedAt          *time.Time
}
