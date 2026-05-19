package usecase

import (
	"fmt"

	appErr "service-core/internal/common/errors"
	"service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/repository"

	"github.com/google/uuid"
)

type RemoveItemUsecase struct {
	cartRepo repository.CartRepository
}

func NewRemoveItemUsecase(
	cartRepo repository.CartRepository,
) *RemoveItemUsecase {
	return &RemoveItemUsecase{
		cartRepo: cartRepo,
	}
}

type RemoveItemInput struct {
	UserID, ProductID, ShopID uuid.UUID
}

func (u *RemoveItemUsecase) Execute(input RemoveItemInput) error {
	if input.ShopID == uuid.Nil {
		return appErr.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(input.UserID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		return appErr.NewNotFound(domain.ErrCartNotFound.Error())
	}

	cart.RemoveItem(input.ProductID, input.ShopID)

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}
