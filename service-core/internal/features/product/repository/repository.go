package repository

import "service-core/internal/features/product/domain"

type FindProductParams struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
}

type ProductRepository interface {
	FindProducts(params FindProductParams) ([]domain.Product, int, error)
	GetById(id string) (domain.Product, error)
}
