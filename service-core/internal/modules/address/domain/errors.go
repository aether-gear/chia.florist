package domain

import "errors"

var (
	ErrAddressLimitReached = errors.New("maximum address limit reached")
)
