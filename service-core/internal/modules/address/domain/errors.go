package domain

import "errors"

var (
	ErrCannotDeleteDefaultAddress = errors.New("default address cannot be deleted")
	ErrAddressNotFound            = errors.New("address not found")
	ErrAddressLimitReached        = errors.New("maximum address limit reached")
)
