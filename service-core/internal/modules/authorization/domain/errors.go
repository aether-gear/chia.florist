package domain

import "errors"

var (
	ErrMerchantRequired    = errors.New("merchant required")
	ErrCustomerRequired    = errors.New("customer required")
	ErrInsufficientRole    = errors.New("insufficient role")
	ErrMemberAlreadyExists = errors.New("member already exists in this merchant")
)

