package domain

import "errors"

var (
	ErrInvalidShopID   = errors.New("invalid shop id")
	ErrInvalidQuantity = errors.New("invalid quantity")

	ErrProductAlreadyAssignedToShop = errors.New("product already exists in cart from another shop")

	ErrInsufficientStock = errors.New("insufficient stock")

	ErrShopCouriersNotFound = errors.New("shop has no active courier service available")

	ErrProductNotFound  = errors.New("product not found")
	ErrCartNotFound     = errors.New("cart not found")
	ErrCartItemNotFound = errors.New("cart item not found")
)
