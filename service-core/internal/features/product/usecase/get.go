package usecase

import "service-core/internal/features/product/domain"

type GetProductUsecase struct {
	repo domain.ProductRepository
}

func NewGetProductsUsecase(repo domain.ProductRepository) *GetProductUsecase {
	return &GetProductUsecase{repo: repo}
}

func (u *GetProductUsecase) Execute(id string) (domain.Product, error) {
	return u.repo.GetById(id)
}
