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

type ProductImageRepository interface {
	Create(images []domain.ProductImage) error

	ListByProductIDs(productIDs []uuid.UUID) (map[uuid.UUID][]domain.ProductImage, error)
	ListByProductID(productID uuid.UUID) ([]domain.ProductImage, error)

	SoftDeleteByProductID(productID uuid.UUID) error
}

type ProductImageUploadService interface {
	Upload(params UploadProductImagesParams) ([]UploadedProductImage, error)
	Delete(assetKeys []string) error
}
