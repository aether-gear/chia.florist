package usecase

import (
	"context"
	"fmt"
	"time"

	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	"service-core/internal/shared/slug"

	"github.com/google/uuid"
)

type CreateShopUsecase struct {
	shopRepo repository.ShopRepository
	slugGen  slug.Generator
}

func NewCreateShopUsecase(
	shopRepo repository.ShopRepository,
	slugGen slug.Generator,
) *CreateShopUsecase {
	return &CreateShopUsecase{
		shopRepo: shopRepo,
		slugGen:  slugGen,
	}
}

type CreateShopInput struct {
	Name        string
	Description *string
	IsActive    bool
}

func (u *CreateShopUsecase) Execute(
	ctx context.Context,
	input CreateShopInput,
) error {
	shop := domain.Shop{
		ID:          uuid.New(),
		Name:        input.Name,
		Slug:        u.slugGen.Generate(input.Name),
		Description: input.Description,
		IsActive:    input.IsActive,
		CreatedAt:   time.Now(),
	}

	err := u.shopRepo.Create(ctx, shop)
	if err != nil {
		return fmt.Errorf("failed to create shop: %w", err)
	}

	return nil
}
