package domain

import "errors"

var (
	ErrNotFoundMerchant = errors.New("merchant not found")
	ErrInvalidName      = errors.New("name is required")
	ErrInvalidEmail     = errors.New("email is required")
)
