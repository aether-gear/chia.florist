package repository

import (
	"errors"
	domain "service-core/internal/features/product/domain"
	"time"
)

var products = []domain.Product{
	{
		ID:          "1",
		Name:        "Rose Bouquet",
		Description: "A beautiful bouquet of fresh red roses",
		Price:       _intPtr(150000),
		Category:    "Flowers",
		Stock:       20,

		Variants: []domain.ProductVariant{},
		Images:   []string{"rose.jpg"},

		Rating:       _float64Ptr(4.5),
		ReviewsCount: _intPtr(120),

		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ArchivedAt: nil,
		DeletedAt:  nil,
	},
	{
		ID:          "2",
		Name:        "Tulip Bundle",
		Description: "Colorful tulip bundle for special occasions",
		Price:       _intPtr(120000),
		Category:    "Flowers",
		Stock:       15,

		Variants: []domain.ProductVariant{},
		Images:   []string{"tulip.jpg"},

		Rating:       _float64Ptr(4.3),
		ReviewsCount: _intPtr(85),

		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ArchivedAt: nil,
		DeletedAt:  nil,
	},
	{
		ID:          "3",
		Name:        "Sunflower Basket",
		Description: "Bright sunflower basket to cheer up your day",
		Price:       _intPtr(100000),
		Category:    "Flowers",
		Stock:       10,

		Variants: []domain.ProductVariant{},
		Images:   []string{"sunflower.jpg"},

		Rating:       _float64Ptr(4.6),
		ReviewsCount: _intPtr(60),

		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		ArchivedAt: nil,
		DeletedAt:  nil,
	},
}

type ProductRepositoryImpl struct{}

func NewProductRepository() *ProductRepositoryImpl {
	return &ProductRepositoryImpl{}
}

func (r *ProductRepositoryImpl) FindProducts(params domain.FindProductParams) ([]domain.Product, int, error) {
	filtered := make([]domain.Product, 0)

	for _, p := range products {
		if p.DeletedAt != nil {
			continue
		}

		if params.ID != nil && p.ID != *params.ID {
			continue
		}

		if params.Name != nil && p.Name != *params.Name {
			continue
		}

		filtered = append(filtered, p)
	}

	total := len(filtered)

	start := (params.Page - 1) * params.Limit
	end := start + params.Limit

	if start > total {
		return []domain.Product{}, total, nil
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *ProductRepositoryImpl) GetById(id string) (domain.Product, error) {
	for _, p := range products {
		if p.DeletedAt != nil {
			continue
		}

		if p.ID == id {
			return p, nil
		}
	}

	return domain.Product{}, errors.New("product not found")
}

func _float64Ptr(v float64) *float64 {
	return &v
}

func _intPtr(v int) *int {
	return &v
}
