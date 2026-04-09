package repository

import domain "service-core/internal/features/product/domain"

type ProductRepositoryImpl struct{}

func NewProductRepository() *ProductRepositoryImpl {
	return &ProductRepositoryImpl{}
}

func (r *ProductRepositoryImpl) FindProducts(params domain.FindProductParams) ([]domain.Product, int, error) {
	products := []domain.Product{
		{
			ProductID: "1",
			Name:      "Rose Bouquet",
			Price:     150000,
			Category:  "Flowers",
			Thumbnail: "rose.jpg",
			Rating:    4.5,
		},
		{
			ProductID: "2",
			Name:      "Tulip Bundle",
			Price:     120000,
			Category:  "Flowers",
			Thumbnail: "tulip.jpg",
			Rating:    4.3,
		},
		{
			ProductID: "3",
			Name:      "Sunflower Basket",
			Price:     100000,
			Category:  "Flowers",
			Thumbnail: "sunflower.jpg",
			Rating:    4.6,
		},
	}

	filtered := make([]domain.Product, 0)

	for _, p := range products {
		if params.ID != nil && p.ProductID != *params.ID {
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
