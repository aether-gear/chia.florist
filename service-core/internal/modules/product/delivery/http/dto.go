package http

import (
	"time"

	"github.com/google/uuid"
)

type productStatusDTO string

const (
	ProductStatusActive   productStatusDTO = "active"
	ProductStatusInactive productStatusDTO = "inactive"
	ProductStatusArchived productStatusDTO = "archived"
)

type createProductRequest struct {
	SKU         string           `json:"sku"`
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	IsAvailable bool             `json:"is_available"`
	Status      productStatusDTO `json:"status"`
	Price       int64            `json:"price"`
	Weight      *float64         `json:"weight"`
}

type productImageResponse struct {
	Thumbnail *string `json:"thumbnail"`
	Preview   *string `json:"preview"`
	Detail    *string `json:"detail"`
}

type productCatalogResponse struct {
	ID           uuid.UUID                     `json:"id"`
	SKU          string                        `json:"sku"`
	Name         string                        `json:"name"`
	Slug         string                        `json:"slug"`
	Status       string                        `json:"status"`
	IsAvailable  bool                          `json:"is_available"`
	Price        int64                         `json:"price"`
	TotalStock   int                           `json:"stock"`
	Image        productImageResponse          `json:"images"`
	Availability []productAvailabilityResponse `json:"availability"`
}

type productAvailabilityResponse struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

type productDetailResponse struct {
	ID          uuid.UUID              `json:"id"`
	SKU         string                 `json:"sku"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description *string                `json:"description"`
	IsAvailable bool                   `json:"is_available"`
	Price       int64                  `json:"price"`
	Weight      *float64               `json:"weight"`
	TotalStock  int                    `json:"stock"`
	UpdatedAt   *time.Time             `json:"updated_at"`
	Images      []productImageResponse `json:"images"`
}

type productInventoryView struct {
	ID         uuid.UUID `json:"id"`
	ShopID     uuid.UUID `json:"shop_id"`
	TotalStock int       `json:"stock"`
	Available  int       `json:"available"`
}
