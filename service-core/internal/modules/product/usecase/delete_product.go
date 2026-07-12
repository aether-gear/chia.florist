package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/product/repository"
	transaction "service-core/internal/shared/transaction"
)

type DeleteProductUsecase struct {
	productRepo repository.ProductRepository
	executor    transaction.Executor
}

func NewDeleteProductUsecase(
	productRepo repository.ProductRepository,
	executor transaction.Executor,
) *DeleteProductUsecase {
	return &DeleteProductUsecase{
		productRepo: productRepo,
		executor:    executor,
	}
}

func (u *DeleteProductUsecase) Execute(
	ctx context.Context,
	slug string,
) error {
	product, err := u.productRepo.
		GetBySlug(ctx, u.executor,
			slug,
		)
	if err != nil {
		return fmt.Errorf("failed to retrieve product: %w", err)
	}

	if product == nil ||
		product.DeletedAt != nil {
		return apperrors.NewNotFound("product not found")
	}

	if err := u.productRepo.
		Delete(ctx, u.executor,
			product.ID,
		); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
