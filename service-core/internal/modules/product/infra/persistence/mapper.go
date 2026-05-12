package persistence

import (
	"service-core/internal/modules/product/domain"
	"service-core/internal/shared/conversion"
)

func (m *productModel) ToDomain() (*domain.Product, error) {
	price, err := conversion.ParsePriceToInt64(m.BasePrice)
	if err != nil {
		return nil, err
	}

	var weight *float64
	if m.Weight != nil {
		w, err := conversion.ParseStringToFloat(m.Weight)
		if err != nil {
			return nil, err
		}
		weight = w
	}

	return &domain.Product{
		ID:          m.ID,
		SKU:         m.SKU,
		Name:        m.Name,
		Description: m.Description,
		Status:      domain.ProductStatus(m.Status),

		Price:  price,
		Weight: weight,

		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		ArchivedAt: m.ArchivedAt,
		DeletedAt:  m.DeletedAt,
	}, nil
}

func FromDomain(p *domain.Product) *productModel {
	basePrice := conversion.FormatInt64ToPrice(p.Price)

	var weight *string
	if p.Weight != nil {
		w := conversion.FormatFloatToString(*p.Weight)
		weight = &w
	}

	return &productModel{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Description: p.Description,
		Status:      string(p.Status),

		BasePrice: basePrice,
		Weight:    weight,

		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		ArchivedAt: p.ArchivedAt,
		DeletedAt:  p.DeletedAt,
	}
}
