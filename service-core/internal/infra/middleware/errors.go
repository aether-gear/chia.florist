package middleware

import "errors"

var (
	ErrAuthenticationRequired = errors.New("authentication required")
)

var (
	ErrInvalidToken = errors.New("invalid access token")
)

var (
	ErrInvalidSession = errors.New("invalid session")
)
