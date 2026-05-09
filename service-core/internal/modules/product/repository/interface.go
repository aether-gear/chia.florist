package repository

import (
	"io"
	"service-core/internal/modules/product/domain"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByID(id uuid.UUID) (*domain.Product, error)

	FindProducts(params FindProductParams) ([]domain.Product, int, error)
	FindByIDs(IDs []uuid.UUID) ([]domain.Product, error)

	CreateProduct(product *domain.Product) error
}

type ImageStorage interface {
	Upload(
		productID uuid.UUID,
		metadata domain.ProductImageMetadata,
		file io.Reader,
	) (*domain.ProductImage, error)
	Delete(objectKey string) error
	Exists(objectKey string) (bool, error)
}
