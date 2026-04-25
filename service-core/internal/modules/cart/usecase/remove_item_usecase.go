package usecase

import (
	"fmt"

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

func (u *RemoveItemUsecase) Execute(userID uuid.UUID, productID uuid.UUID) error {
	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}
	if cart == nil {
		return fmt.Errorf("failed to load cart with items: cart not found")
	}

	cart.RemoveItem(productID)

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}
