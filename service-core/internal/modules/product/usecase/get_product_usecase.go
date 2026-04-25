package usecase

import (
	"fmt"

	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetProductUsecase struct {
	repo repository.ProductRepository
}

func NewGetProductsUsecase(
	repo repository.ProductRepository,
) *GetProductUsecase {
	return &GetProductUsecase{
		repo: repo,
	}
}

func (u *GetProductUsecase) Execute(id uuid.UUID) (*repository.ProductWithInventory, error) {
	product, err := u.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load products with inventory: %w", err)
	}

	return product, nil
}
