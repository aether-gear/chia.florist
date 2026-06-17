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
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Product, error)

	FindProducts(
		ctx context.Context,
		exec transaction.Executor,
		params FindProductParams,
	) ([]domain.Product, int, error)

	FindProductsWithInventory(
		ctx context.Context,
		exec transaction.Executor,
		params FindProductParams,
	) ([]domain.ProductWithInventory, int, error)

	FindByIDs(
		ctx context.Context,
		exec transaction.Executor,
		IDs []uuid.UUID,
	) ([]domain.Product, error)

	CreateProduct(
		ctx context.Context,
		exec transaction.Executor,
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
		exec transaction.Executor,
		productIDs []uuid.UUID,
	) (map[uuid.UUID][]domain.ProductImage, error)

	ListByProductID(
		ctx context.Context,
		exec transaction.Executor,
		productID uuid.UUID,
	) ([]domain.ProductImage, error)

	SoftDeleteByProductID(
		ctx context.Context,
		exec transaction.Executor,
		productID uuid.UUID,
	) error
}

type ProductImageUploadService interface {
	Upload(params UploadProductImagesParams) ([]UploadedProductImage, error)
	Delete(assetKeys []string) error
}
