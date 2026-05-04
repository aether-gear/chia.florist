package repository

import (
	"service-core/internal/modules/product/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByID(id uuid.UUID) (*domain.Product, error)

	FindProducts(params FindProductParams) ([]domain.Product, int, error)
	FindByIDs(IDs []uuid.UUID) ([]domain.Product, error)

	CreateProduct(product *domain.Product) error
}
