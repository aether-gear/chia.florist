package usecase

import (
	"fmt"

	cartR "service-core/internal/modules/cart/repository"
	productR "service-core/internal/modules/product/repository"

	"github.com/google/uuid"
)

type AddItemUsecase struct {
	cartRepo    cartR.CartRepository
	productRepo productR.ProductRepository
}

func NewAddItemUsecase(
	cR cartR.CartRepository,
	pR productR.ProductRepository,
) *AddItemUsecase {
	return &AddItemUsecase{
		cartRepo:    cR,
		productRepo: pR,
	}
}

func (u *AddItemUsecase) Execute(userID uuid.UUID, productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("failed to add item: invalid quantity")
	}

	product, err := u.productRepo.GetByID(productID)
	if err != nil {
		return fmt.Errorf("failed to load product with inventory: %w", err)
	}
	if product == nil {
		return fmt.Errorf("failed to load product with inventory: product not found")
	}

	if quantity > product.Inventory.Stock {
		return fmt.Errorf("failed to add item: insufficient stock")
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

	if err := cart.AddItem(productID, quantity); err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	if err := u.cartRepo.Save(cart); err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	return nil
}
