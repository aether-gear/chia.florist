package http

import (
	"service-core/internal/features/product/domain"
	"service-core/internal/features/product/repository"
	"time"

	"github.com/google/uuid"
)

type ProductOverviewResponse struct {
	ID     uuid.UUID            `json:"id"`
	SKU    string               `json:"sku"`
	Name   string               `json:"name"`
	Status domain.ProductStatus `json:"status"`

	Price int64 `json:"price"`

	Stock         int
	ReservedStock int
}

type ProductDetailResponse struct {
	ID          uuid.UUID            `json:"id"`
	SKU         string               `json:"sku"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Status      domain.ProductStatus `json:"status"`

	Price  int64    `json:"price"`
	Weight *float64 `json:"weight"`

	Stock         int
	ReservedStock int

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

func ToListResponse(p repository.ProductWithInventory) ProductOverviewResponse {
	return ProductOverviewResponse{
		ID:            p.Product.ID,
		SKU:           p.Product.SKU,
		Name:          p.Product.Name,
		Status:        p.Product.Status,
		Price:         p.Product.Price,
		Stock:         p.Inventory.Stock,
		ReservedStock: p.Inventory.ReservedStock,
	}
}

func ToDetailResponse(p repository.ProductWithInventory) ProductDetailResponse {
	return ProductDetailResponse{
		ID:            p.Product.ID,
		SKU:           p.Product.SKU,
		Name:          p.Product.Name,
		Description:   p.Product.Description,
		Status:        p.Product.Status,
		Price:         p.Product.Price,
		Weight:        p.Product.Weight,
		Stock:         p.Inventory.Stock,
		ReservedStock: p.Inventory.ReservedStock,
		CreatedAt:     p.Product.CreatedAt,
		UpdatedAt:     p.Product.UpdatedAt,
		ArchivedAt:    p.Product.ArchivedAt,
	}
}
