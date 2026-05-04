package http

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatusDTO string

const (
	ProductStatusActive   ProductStatusDTO = "active"
	ProductStatusInactive ProductStatusDTO = "inactive"
	ProductStatusArchived ProductStatusDTO = "archived"
)

type ProductOverviewResponse struct {
	ID            uuid.UUID        `json:"id"`
	SKU           string           `json:"sku"`
	Name          string           `json:"name"`
	Slug          string           `json:"slug"`
	Status        ProductStatusDTO `json:"status"`
	Price         int64            `json:"price"`
	Stock         int              `json:"stock"`
	ReservedStock int              `json:"reserved_stock"`
}

type ProductDetailResponse struct {
	ID            uuid.UUID              `json:"id"`
	SKU           string                 `json:"sku"`
	Name          string                 `json:"name"`
	Slug          string                 `json:"slug"`
	Description   *string                `json:"description"`
	Status        ProductStatusDTO       `json:"status"`
	Price         int64                  `json:"price"`
	Weight        *float64               `json:"weight"`
	Stock         int                    `json:"stock"`
	ReservedStock int                    `json:"reserved_stock"`
	Inventories   []ProductInventoryView `json:"inventories"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     *time.Time             `json:"updated_at"`
	ArchivedAt    *time.Time             `json:"archived_at"`
}

type ProductInventoryView struct {
	ID        uuid.UUID `json:"id"`
	ShopID    uuid.UUID `json:"shop_id"`
	Stock     int       `json:"stock"`
	Reserved  int       `json:"reserved"`
	Available int       `json:"available"`
}

type CreateProductRequest struct {
	SKU         string           `json:"sku"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	Status      ProductStatusDTO `json:"status"`
	Price       int64            `json:"price"`
	Weight      *float64         `json:"weight"`
}
