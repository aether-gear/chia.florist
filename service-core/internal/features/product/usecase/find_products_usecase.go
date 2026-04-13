package usecase

import (
	"service-core/internal/features/product/repository"
)

type FindProductsUsecase struct {
	repo repository.ProductRepository
}

func NewFindProductsUsecase(repo repository.ProductRepository) *FindProductsUsecase {
	return &FindProductsUsecase{repo: repo}
}

func (u *FindProductsUsecase) Execute(params repository.FindProductParams) ([]repository.ProductWithInventory, int, error) {
	return u.repo.FindProducts(params)
}
