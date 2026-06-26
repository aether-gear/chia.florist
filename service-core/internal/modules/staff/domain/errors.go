package domain

import "errors"

var (
	ErrNotFoundStaff = errors.New("staff not found")
	// ErrInvalidName      = errors.New("name is required")
	ErrInvalidEmail = errors.New("email is required")
)
