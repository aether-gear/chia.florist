package usecase

import (
	"fmt"
	"service-core/internal/modules/product/domain"
	"service-core/internal/modules/product/repository"
	"time"

	"github.com/google/uuid"
)

type CreateProductUsecase struct {
	productRepo repository.ProductRepository
}

func NewCreateProductUsecase(pR repository.ProductRepository) *CreateProductUsecase {
	return &CreateProductUsecase{
		productRepo: pR,
	}
}

type CreateProductInput struct {
	SKU          string
	Name         string
	Description  *string
	Status       domain.ProductStatus
	Price        int64
	Weight       *float64
	InitialStock int
}

func (u *CreateProductUsecase) Execute(input CreateProductInput) error {
	now := time.Now()

	product := &domain.Product{
		ID:          uuid.New(),
		SKU:         input.SKU,
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		Price:       input.Price,
		Weight:      input.Weight,
		CreatedAt:   now,
	}
	if err := product.Validate(); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	inventory := &domain.Inventory{
		ID:            uuid.New(),
		ProductID:     product.ID,
		Stock:         input.InitialStock,
		ReservedStock: 0,
		CreatedAt:     now,
	}
	if err := inventory.Validate(); err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}

	if err := u.productRepo.CreateProduct(product); err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}

	if err := u.productRepo.CreateInventory(inventory); err != nil {
		return fmt.Errorf("failed to save inventory: %w", err)
	}

	return nil
}
