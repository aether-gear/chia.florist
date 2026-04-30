package domain

import "errors"

var (
	ErrInvalidCourier     = errors.New("courier is required")
	ErrInvalidService     = errors.New("service is required")
	ErrInvalidOrigin      = errors.New("origin is required")
	ErrInvalidDestination = errors.New("destination is required")

	ErrInvalidCost   = errors.New("cost must be greater than 0")
	ErrInvalidWeight = errors.New("weight must be greater than 0")

	ErrInvalidStatus = errors.New("invalid status")
)

var (
	ErrInvalidEventStatus      = ErrInvalidStatus
	ErrInvalidEventDescription = errors.New("event description is required")
)

var (
	ErrInvalidDisplayName = errors.New("display name is required")
)
