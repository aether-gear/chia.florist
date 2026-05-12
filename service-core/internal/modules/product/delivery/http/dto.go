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

type ProductImageResponse struct {
	Thumbnail *string `json:"thumbnail,omitempty"`
	Preview   *string `json:"preview,omitempty"`
	Detail    *string `json:"detail,omitempty"`
}

type ProductCatalogResponse struct {
	ID          uuid.UUID            `json:"id"`
	SKU         string               `json:"sku"`
	Name        string               `json:"name"`
	Slug        string               `json:"slug"`
	IsAvailable bool                 `json:"is_available"`
	Price       int64                `json:"price"`
	TotalStock  int                  `json:"stock"`
	Image       ProductImageResponse `json:"images"`
}

type ProductDetailResponse struct {
	ID          uuid.UUID  `json:"id"`
	SKU         string     `json:"sku"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description"`
	IsAvailable bool       `json:"is_available"`
	Price       int64      `json:"price"`
	Weight      *float64   `json:"weight"`
	TotalStock  int        `json:"stock"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type ProductInventoryView struct {
	ID         uuid.UUID `json:"id"`
	ShopID     uuid.UUID `json:"shop_id"`
	TotalStock int       `json:"stock"`
	Available  int       `json:"available"`
}

type CreateProductRequest struct {
	SKU         string           `json:"sku"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	IsAvailable bool             `json:"is_available"`
	Status      ProductStatusDTO `json:"status"`
	Price       int64            `json:"price"`
	Weight      *float64         `json:"weight"`
}
