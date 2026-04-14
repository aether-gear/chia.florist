package usecase

import (
	"fmt"
	"service-core/internal/features/product/domain"
	"service-core/internal/features/product/repository"
	"time"

	"github.com/google/uuid"
)

type CreateProductInput struct {
	SKU          string
	Name         string
	Description  *string
	Status       domain.ProductStatus
	Price        int64
	Weight       *float64
	InitialStock int
}

type CreateProductUsecase struct {
	productRepo repository.ProductRepository
}

func NewCreateProductUsecase(pR repository.ProductRepository) *CreateProductUsecase {
	return &CreateProductUsecase{
		productRepo: pR,
	}
}

func (u *CreateProductUsecase) Execute(input CreateProductInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if input.InitialStock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}

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

	inventory := &domain.Inventory{
		ID:            uuid.New(),
		ProductID:     product.ID,
		Stock:         input.InitialStock,
		ReservedStock: 0,
		CreatedAt:     now,
	}

	if err := u.productRepo.CreateProduct(product); err != nil {
		return err
	}

	if err := u.productRepo.CreateInventory(inventory); err != nil {
		return err
	}

	return nil
}
