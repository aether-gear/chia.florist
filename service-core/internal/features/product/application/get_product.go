package application

import "service-core/internal/features/product/domain"

type GetProductUsecase struct {
	repo domain.ProductRepository
}

func NewGetProductsUsecase(repo domain.ProductRepository) *GetProductUsecase {
	return &GetProductUsecase{repo: repo}
}

func (u *GetProductUsecase) Execute(id string) (*domain.Product, error) {
	query := domain.FindProductParams{
		ID:    &id,
		Page:  1,
		Limit: 1,
	}

	products, _, err := u.repo.FindProducts(query)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	return &products[0], nil
}
