package usecase

import (
	"fmt"
	"service-core/internal/features/cart/repository"

	"github.com/google/uuid"
)

type RemoveItemUsecase struct {
	cartRepo repository.CartRepository
}

func NewRemoveItemUsecase(cR repository.CartRepository) *RemoveItemUsecase {
	return &RemoveItemUsecase{
		cartRepo: cR,
	}
}

func (u *RemoveItemUsecase) Execute(userID uuid.UUID, productID uuid.UUID) error {
	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return err
	}
	if cart == nil {
		return fmt.Errorf("cart not found")
	}

	cart.RemoveItem(productID)

	return u.cartRepo.Save(cart)
}
