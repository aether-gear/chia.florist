package application

import "service-core/internal/features/product/domain"

type FindProductsUsecase struct {
	repo domain.ProductRepository
}

func NewFindProductsUsecase(repo domain.ProductRepository) *FindProductsUsecase {
	return &FindProductsUsecase{repo: repo}
}

func (u *FindProductsUsecase) Execute(params domain.FindProductParams) ([]domain.Product, int, error) {
	return u.repo.FindProducts(params)
}
