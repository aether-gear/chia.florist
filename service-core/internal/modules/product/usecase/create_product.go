package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	"service-core/internal/shared/slug"

	"github.com/google/uuid"
)

type CreateProductUsecase struct {
	productRepo repository.ProductRepository
	slugGen     slug.Generator
}

func NewCreateProductUsecase(
	productRepo repository.ProductRepository,
	slugGen slug.Generator,
) *CreateProductUsecase {
	return &CreateProductUsecase{
		productRepo: productRepo,
		slugGen:     slugGen,
	}
}

type CreateProductInput struct {
	SKU         string
	Name        string
	Description *string
	Status      string
	Price       int64
	Weight      *float64
}

func (u *CreateProductUsecase) Execute(
	ctx context.Context,
	input CreateProductInput,
) error {
	now := time.Now()

	product := &domain.Product{
		ID:          uuid.New(),
		SKU:         input.SKU,
		Name:        input.Name,
		Slug:        u.slugGen.Generate(input.Name),
		Description: input.Description,
		Status:      domain.ProductStatus(input.Status),
		Price:       input.Price,
		Weight:      input.Weight,
		CreatedAt:   now,
	}
	if err := product.Validate(); err != nil {
		if errors.Is(err, domain.ErrInvalidProductName) ||
			errors.Is(err, domain.ErrInvalidProductPrice) {
			return apperrors.NewInvalidInput(err.Error())
		}

		return err
	}

	if err := u.productRepo.CreateProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}

	return nil
}
