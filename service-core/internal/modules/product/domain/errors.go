package domain

import "errors"

var ErrInvalidStock = errors.New("stock cannot be negative")
var ErrInvalidProductName = errors.New("name is required")
var ErrInvalidProductPrice = errors.New("price must be greater than 0")
