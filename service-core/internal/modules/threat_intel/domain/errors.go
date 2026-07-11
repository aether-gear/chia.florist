package domain

import "errors"

var (
	ErrProviderUnavailable = errors.New("threat intel provider unavailable")
	ErrInvalidIP           = errors.New("invalid ip address")
	ErrAPIKeyRequired      = errors.New("api key is required but not configured")
)
