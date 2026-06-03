package repository

import (
	"context"

	"service-core/internal/modules/product/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ProductRepository interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Product, error)

	FindProducts(
		ctx context.Context,
		params FindProductParams,
	) ([]domain.Product, int, error)
	FindByIDs(
		ctx context.Context,
		IDs []uuid.UUID,
	) ([]domain.Product, error)

	CreateProduct(
		ctx context.Context,
		product *domain.Product,
	) error
}

type ProductImageRepository interface {
	Create(
		ctx context.Context,
		exec transaction.Executor,
		images []domain.ProductImage,
	) error

	ListByProductIDs(
		ctx context.Context,
		productIDs []uuid.UUID,
	) (map[uuid.UUID][]domain.ProductImage, error)
	ListByProductID(
		ctx context.Context,
		productID uuid.UUID,
	) ([]domain.ProductImage, error)

	SoftDeleteByProductID(
		ctx context.Context,
		productID uuid.UUID,
	) error
}

type ProductImageUploadService interface {
	Upload(params UploadProductImagesParams) ([]UploadedProductImage, error)
	Delete(assetKeys []string) error
}
