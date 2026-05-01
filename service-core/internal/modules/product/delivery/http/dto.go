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
	ID     uuid.UUID        `json:"id"`
	SKU    string           `json:"sku"`
	Name   string           `json:"name"`
	Slug   string           `json:"slug"`
	Status ProductStatusDTO `json:"status"`

	Price int64 `json:"price"`

	Stock         int
	ReservedStock int
}

type ProductDetailResponse struct {
	ID          uuid.UUID        `json:"id"`
	SKU         string           `json:"sku"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description *string          `json:"description"`
	Status      ProductStatusDTO `json:"status"`

	Price  int64    `json:"price"`
	Weight *float64 `json:"weight"`

	Stock         int
	ReservedStock int

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

type CreateProductRequest struct {
	SKU          string           `json:"sku"`
	Name         string           `json:"name"`
	Description  *string          `json:"description"`
	Status       ProductStatusDTO `json:"status"`
	Price        int64            `json:"price"`
	Weight       *float64         `json:"weight"`
	InitialStock int              `json:"stock"`
}
