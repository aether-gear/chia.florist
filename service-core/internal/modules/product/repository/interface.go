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
	CreateProductImage(image *domain.ProductImage) error
	CreateDetailImages(images []domain.ProductDetailImage) error

	GetByProductID(productID uuid.UUID) ([]domain.ProductImage, error)

	DeleteByProductID(productID uuid.UUID) error
}

type ProductImageUploadRepository interface {
	UploadCatalogImage(params UploadProductImageParams) (string, error)
	UploadCartImage(params UploadProductImageParams) (string, error)
	UploadDetailImage(params UploadProductImageParams) (string, error)

	DeleteUploadedImages(urls []string) error
}
