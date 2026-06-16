package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"
)

type ListShopsUsecase struct {
	executor transaction.Executor
	shopRepo repository.ShopRepository
}

func NewListShopsUsecase(
	executor transaction.Executor,
	shopRepo repository.ShopRepository,
) *ListShopsUsecase {
	return &ListShopsUsecase{
		executor: executor,
		shopRepo: shopRepo,
	}
}

type ListShopsInput struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
}

func (u *ListShopsUsecase) Execute(
	ctx context.Context,
	input ListShopsInput,
) ([]domain.Shop, int, error) {
	params := repository.FindShopsParams{
		Page:  input.Page,
		Limit: input.Limit,
		ID:    input.ID,
		Name:  input.Name,
	}

	shops, total, err := u.shopRepo.List(ctx, u.executor, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load shops: %w", err)
	}
	if len(shops) == 0 {
		return nil, 0, apperrors.NewNotFound("shops not available at the moment")
	}

	return shops, total, nil
}
