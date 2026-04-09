package http

import (
	"service-core/internal/features/product/domain"
	"time"
)

type ProductDetailResponse struct {
	ID          string
	Name        string
	Description string
	Price       *int
	Category    string
	Stock       int

	Variants []domain.ProductVariant
	Images   []string

	Rating       *float64
	ReviewsCount *int

	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time
}

type ProductOverviewResponse struct {
	ID       string
	Name     string
	Price    *int
	Category string
	Stock    int

	Rating       *float64
	ReviewsCount *int
}

func ToListResponse(p domain.Product) ProductOverviewResponse {
	return ProductOverviewResponse{
		ID:           p.ID,
		Name:         p.Name,
		Price:        p.Price,
		Category:     p.Category,
		Stock:        p.Stock,
		Rating:       p.Rating,
		ReviewsCount: p.ReviewsCount,
	}
}

func ToDetailResponse(p domain.Product) ProductDetailResponse {
	return ProductDetailResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Stock:       p.Stock,

		Variants: p.Variants,
		Images:   p.Images,

		Rating:       p.Rating,
		ReviewsCount: p.ReviewsCount,

		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		ArchivedAt: p.ArchivedAt,
	}
}
