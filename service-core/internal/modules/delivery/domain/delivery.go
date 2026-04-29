package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus string

const (
	Pending   ShipmentStatus = "pending"
	Paid      ShipmentStatus = "paid"
	Shipped   ShipmentStatus = "shipped"
	InTransit ShipmentStatus = "in_transit"
	Delivered ShipmentStatus = "delivered"
	Cancelled ShipmentStatus = "cancelled"
)

type Delivery struct {
	ID uuid.UUID

	OrderID  uuid.UUID
	MethodID uuid.UUID

	Status         ShipmentStatus
	TrackingNumber *string

	Courier string
	Service string

	Cost          int64
	Weight        int
	OriginID      string
	DestinationID string

	CreatedAt time.Time
}

func (d *Delivery) Validate() error {
	if d.Cost <= 0 {
		return ErrInvalidCost
	}

	if d.Weight <= 0 {
		return ErrInvalidWeight
	}

	if d.Courier == "" {
		return ErrInvalidCourier
	}

	if d.Service == "" {
		return ErrInvalidService
	}

	if d.OriginID == "" {
		return ErrInvalidOrigin
	}

	if d.DestinationID == "" {
		return ErrInvalidDestination
	}

	return nil
}

func (s ShipmentStatus) IsValid() bool {
	switch s {
	case Pending, Paid, Shipped, InTransit, Delivered, Cancelled:
		return true
	default:
		return false
	}
}
