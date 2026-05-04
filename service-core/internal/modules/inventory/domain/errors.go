package domain

import "errors"

var (
	ErrInvalidStock    = errors.New("stock cannot be negative")
	ErrInvalidReserved = errors.New("reserved stock cannot be negative")
)
