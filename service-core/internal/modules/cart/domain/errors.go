package domain

import "errors"

var (
	ErrInvalidQuantity = errors.New("invalid quantity")

	ErrInsufficientStock = errors.New("insufficient stock")

	ErrProductNotFound  = errors.New("product not found")
	ErrCartNotFound     = errors.New("cart not found")
	ErrCartItemNotFound = errors.New("cart item not found")
)
