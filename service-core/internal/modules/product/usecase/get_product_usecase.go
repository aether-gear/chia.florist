package usecase

import (
	"service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type GetProductUsecase struct {
	repo repository.ProductRepository
}

func NewGetProductsUsecase(repo repository.ProductRepository) *GetProductUsecase {
	return &GetProductUsecase{repo: repo}
}

func (u *GetProductUsecase) Execute(id uuid.UUID) (*repository.ProductWithInventory, error) {
	return u.repo.GetByID(id)
}
