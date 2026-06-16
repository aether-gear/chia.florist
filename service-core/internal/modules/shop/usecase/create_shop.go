package usecase

import (
	"context"
	"fmt"
	"time"

	authorDomain "service-core/internal/modules/authorization/domain"
	"service-core/internal/modules/shop/domain"
	"service-core/internal/modules/shop/repository"
	slug "service-core/internal/shared/slug"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type CreateShopUsecase struct {
	shopRepo repository.ShopRepository
	slugGen  slug.Generator
	executor transaction.Executor
}

func NewCreateShopUsecase(
	shopRepo repository.ShopRepository,
	slugGen slug.Generator,
	executor transaction.Executor,
) *CreateShopUsecase {
	return &CreateShopUsecase{
		shopRepo: shopRepo,
		slugGen:  slugGen,
		executor: executor,
	}
}

type CreateShopInput struct {
	Name        string
	Description *string
	IsActive    bool
}

func (u *CreateShopUsecase) Execute(
	ctx context.Context,
	actor authorDomain.Actor,
	input CreateShopInput,
) error {
	canSetActive := false
	for _, actorRole := range actor.Roles {
		if actorRole.Code == authorDomain.RoleMerchantAdmin {
			canSetActive = true
			break
		}
	}

	shop := domain.Shop{
		ID:          uuid.New(),
		Name:        input.Name,
		Slug:        u.slugGen.Generate(input.Name),
		Description: input.Description,
		CreatedAt:   time.Now(),
	}
	if canSetActive {
		shop.IsActive = input.IsActive
	}

	err := u.shopRepo.Create(ctx, u.executor, shop)
	if err != nil {
		return fmt.Errorf("failed to create shop: %w", err)
	}

	return nil
}
