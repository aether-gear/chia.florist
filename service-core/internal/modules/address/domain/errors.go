package domain

import "errors"

var (
	ErrCannotDeleteDefaultAddress = errors.New("default address cannot be deleted")
	ErrNotFoundDefaultAddress     = errors.New("default address not found")
	ErrAddressNotFound            = errors.New("address not found")
	ErrAddressLimitReached        = errors.New("maximum address limit reached")
)
