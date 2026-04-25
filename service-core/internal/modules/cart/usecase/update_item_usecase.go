package usecase

import (
	"fmt"

	cartR "service-core/internal/modules/cart/repository"
	productR "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type UpdateItemUsecase struct {
	cartRepo    cartR.CartRepository
	productRepo productR.ProductRepository
}

func NewUpdateItemUsecase(
	cR cartR.CartRepository,
	pR productR.ProductRepository,
) *UpdateItemUsecase {
	return &UpdateItemUsecase{
		cartRepo:    cR,
		productRepo: pR,
	}
}

func (u *UpdateItemUsecase) Execute(userID uuid.UUID, productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("failed to update cart item: invalid quantity")
	}

	product, err := u.productRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("failed to retrieve product: %w", err)
	}
	if product == nil {
		return fmt.Errorf("failed to retrieve product: product not found")
	}

	if quantity > product.Inventory.Stock {
		return fmt.Errorf("failed to update cart item: insufficient stock")
	}

	cart, err := u.cartRepo.GetWithItemsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to load cart with items: %w", err)
	}

	if cart == nil {
		cart, err = u.cartRepo.NewCart(userID)
		if err != nil {
			return fmt.Errorf("failed to create cart: %w", err)
		}
	}

	if !cart.HasItem(productID) {
		return fmt.Errorf("failed to update cart item: item not found")
	}

	if err := cart.SetItem(productID, quantity); err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	return nil
}
