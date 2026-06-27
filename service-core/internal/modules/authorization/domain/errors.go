package domain

import "errors"

var (
	ErrStaffRequired       = errors.New("staff required")
	ErrCustomerRequired    = errors.New("customer required")
	ErrInsufficientRole    = errors.New("insufficient role")
	ErrMemberAlreadyExists = errors.New("member already exists in this staff")
)
