package usecase

import (
	"fmt"

	"service-core/internal/modules/product/repository"
)

type FindProductsUsecase struct {
	repo repository.ProductRepository
}

func NewFindProductsUsecase(
	repo repository.ProductRepository,
) *FindProductsUsecase {
	return &FindProductsUsecase{
		repo: repo,
	}
}

func (u *FindProductsUsecase) Execute(
	params repository.FindProductParams,
) (
	[]repository.ProductWithInventory,
	int,
	error,
) {
	products, total, err := u.repo.FindProducts(params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load products with inventory: %w", err)
	}

	return products, total, nil
}
