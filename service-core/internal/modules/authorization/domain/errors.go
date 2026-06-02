package domain

import "errors"

var (
	ErrMerchantRequired = errors.New("merchant required")
	ErrCustomerRequired = errors.New("customer required")
)
