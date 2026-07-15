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
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type SaveProductUsecase struct {
	productRepo repository.ProductRepository
	slugGen     slug.Generator
	executor    transaction.Executor
	perfRepo    repository.ProductPerformanceRepository
}

func NewSaveProductUsecase(
	productRepo repository.ProductRepository,
	slugGen slug.Generator,
	executor transaction.Executor,
	perfRepo repository.ProductPerformanceRepository,
) *SaveProductUsecase {
	return &SaveProductUsecase{
		productRepo: productRepo,
		slugGen:     slugGen,
		executor:    executor,
		perfRepo:    perfRepo,
	}
}

type SaveProductInput struct {
	ID                   *uuid.UUID
	SKU                  string
	Name                 string
	Description          *string
	Status               string
	Price                int64
	Weight               *float64
	CostPrice            *int64
	SupplierLeadTimeDays *int
}

func (u *SaveProductUsecase) Execute(
	ctx context.Context,
	input SaveProductInput,
) error {
	now := time.Now()

	var productID uuid.UUID
	if input.ID == nil {
		productID = uuid.New()
	} else {
		productID = *input.ID
	}

	product := &domain.Product{
		ID:          productID,
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

	if err := u.productRepo.
		Save(ctx, u.executor,
			product,
		); err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}

	perf := domain.ProductPerformance{
		ProductID:            productID,
		CostPrice:            input.CostPrice,
		SupplierLeadTimeDays: input.SupplierLeadTimeDays,
	}
	if err := u.perfRepo.UpsertPerformance(ctx, u.executor, perf, product.Price); err != nil {
		return fmt.Errorf("failed to save product performance: %w", err)
	}

	return nil
}
