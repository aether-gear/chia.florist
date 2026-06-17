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

type SaveShopUsecase struct {
	shopRepo repository.ShopRepository
	slugGen  slug.Generator
	executor transaction.Executor
}

func NewSaveShopUsecase(
	shopRepo repository.ShopRepository,
	slugGen slug.Generator,
	executor transaction.Executor,
) *SaveShopUsecase {
	return &SaveShopUsecase{
		shopRepo: shopRepo,
		slugGen:  slugGen,
		executor: executor,
	}
}

type SaveShopInput struct {
	ID          *uuid.UUID
	Name        string
	Description *string
	IsActive    bool
}

func (u *SaveShopUsecase) Execute(
	ctx context.Context,
	actor authorDomain.Actor,
	input SaveShopInput,
) error {
	canSetActive := false
	for _, actorRole := range actor.Roles {
		if actorRole.Code == authorDomain.RoleMerchantAdmin {
			canSetActive = true
			break
		}
	}

	var shopID uuid.UUID
	isCreate := input.ID == nil
	if isCreate {
		shopID = uuid.New()
	} else {
		shopID = *input.ID
	}

	shop := domain.Shop{
		ID:          shopID,
		Name:        input.Name,
		Slug:        u.slugGen.Generate(input.Name),
		Description: input.Description,
		CreatedAt:   time.Now(),
	}
	if canSetActive {
		shop.IsActive = input.IsActive
	}

	err := u.shopRepo.
		Save(ctx, u.executor, shop)
	if err != nil {
		return fmt.Errorf("failed to create shop: %w", err)
	}

	return nil
}
