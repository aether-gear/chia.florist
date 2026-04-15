package usecase

import (
	"fmt"
	cartR "service-core/internal/features/cart/repository"
	productR "service-core/internal/features/product/repository"

	"github.com/google/uuid"
)

type AddItemUsecase struct {
	cartRepo    cartR.CartRepository
	productRepo productR.ProductRepository
}

func NewAddItemUsecase(cR cartR.CartRepository, pR productR.ProductRepository) *AddItemUsecase {
	return &AddItemUsecase{
		cartRepo:    cR,
		productRepo: pR,
	}
}

func (u *AddItemUsecase) Execute(userID uuid.UUID, productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("invalid quantity")
	}

	product, err := u.productRepo.GetByID(productID)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("product not found")
	}

	if quantity > product.Inventory.Stock {
		return fmt.Errorf("insufficient stock")
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return err
	}
	if cart == nil {
		cart, err = u.cartRepo.NewCart(userID)
		if err != nil {
			return err
		}
	}

	if err := cart.AddItem(productID, quantity); err != nil {
		return err
	}

	return u.cartRepo.Save(cart)
}
