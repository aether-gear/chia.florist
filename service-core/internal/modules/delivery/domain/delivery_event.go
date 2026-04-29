package domain

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryEvent struct {
	ID uuid.UUID

	ShipmentID uuid.UUID

	Status      string
	Description string
	Location    string

	Timestamp time.Time
}

func (dE *DeliveryEvent) Validate() error {
	if dE.Status == "" {
		return ErrInvalidEventStatus
	}

	if dE.Description == "" {
		return ErrInvalidEventDescription
	}

	return nil
}
