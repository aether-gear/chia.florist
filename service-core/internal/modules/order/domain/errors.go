package domain

import "errors"

var (
	ErrMissingSLAFields = errors.New("orders in confirmed or processing status must have confirmed_at and handling_expires_at populated")
)
