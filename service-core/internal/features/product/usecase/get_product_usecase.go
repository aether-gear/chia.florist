package usecase

import (
	"service-core/internal/features/product/domain"
	"service-core/internal/features/product/repository"
)

type GetProductUsecase struct {
	repo repository.ProductRepository
}

func NewGetProductsUsecase(repo repository.ProductRepository) *GetProductUsecase {
	return &GetProductUsecase{repo: repo}
}

func (u *GetProductUsecase) Execute(id string) (*domain.Product, error) {
	return u.repo.GetById(id)
}
