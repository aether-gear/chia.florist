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
	cR repository.CartRepository,
) *RemoveItemUsecase {
	return &RemoveItemUsecase{
		cartRepo: cR,
	}
}

func (u *RemoveItemUsecase) Execute(userID uuid.UUID, productID uuid.UUID, shopID uuid.UUID) error {
	if shopID == uuid.Nil {
		return appErr.NewInvalidInput(domain.ErrInvalidShopID.Error())
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		return appErr.NewNotFound(domain.ErrCartNotFound.Error())
	}

	cart.RemoveItem(productID, shopID)

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}
