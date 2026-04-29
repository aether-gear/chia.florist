package domain

import "github.com/google/uuid"

type DeliveryOption struct {
	ID uuid.UUID

	IsActive bool

	Courier     string
	Service     string
	DisplayName string

	Description string
}

func (dO *DeliveryOption) Validate() error {
	if dO.Courier == "" {
		return ErrInvalidCourier
	}

	if dO.Service == "" {
		return ErrInvalidService
	}

	if dO.DisplayName == "" {
		return ErrInvalidDisplayName
	}

	return nil
}
