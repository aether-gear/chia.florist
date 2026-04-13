package http

import (
	"service-core/internal/features/product/domain"
	"time"

	"github.com/google/uuid"
)

type SignUpParams struct {
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
}

type SignInEmailParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUsernameParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ProductOverviewResponse struct {
	ID     uuid.UUID            `json:"id"`
	SKU    string               `json:"sku"`
	Name   string               `json:"name"`
	Status domain.ProductStatus `json:"status"`

	Price int64 `json:"price"`
}

type ProductDetailResponse struct {
	ID          uuid.UUID            `json:"id"`
	SKU         string               `json:"sku"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Status      domain.ProductStatus `json:"status"`

	Price  int64    `json:"price"`
	Weight *float64 `json:"weight"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	ArchivedAt *time.Time `json:"archived_at"`
}

func ToListResponse(p domain.Product) ProductOverviewResponse {
	return ProductOverviewResponse{
		ID:     p.ID,
		SKU:    p.SKU,
		Name:   p.Name,
		Status: p.Status,
		Price:  p.Price,
	}
}

func ToDetailResponse(p domain.Product) ProductDetailResponse {
	return ProductDetailResponse{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status,
		Price:       p.Price,
		Weight:      p.Weight,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		ArchivedAt:  p.ArchivedAt,
	}
}
