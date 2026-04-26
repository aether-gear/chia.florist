package domain

import "errors"

var ErrInvalidQuantity = errors.New("invalid quantity")

var ErrInsufficient = errors.New("insufficient stock")

var ErrProductNotFound = errors.New("product not found")
var ErrCartNotFound = errors.New("cart not found")
var ErrCartItemNotFound = errors.New("cart item not found")
