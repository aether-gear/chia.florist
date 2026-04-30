package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentEvent struct {
	ID uuid.UUID

	ShipmentID uuid.UUID

	Status      string
	Description string
	Location    string

	Timestamp time.Time
}

func (dE *ShipmentEvent) Validate() error {
	if dE.Status == "" {
		return ErrInvalidEventStatus
	}

	if dE.Description == "" {
		return ErrInvalidEventDescription
	}

	return nil
}
