package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShipmentStatus string

const (
	ShipmentStatusCreated        ShipmentStatus = "created"
	ShipmentStatusPacked         ShipmentStatus = "packed"
	ShipmentStatusLabelled       ShipmentStatus = "labelled"
	ShipmentStatusPickedUp       ShipmentStatus = "picked_up"
	ShipmentStatusInTransit      ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered      ShipmentStatus = "delivered"
	ShipmentStatusFailed         ShipmentStatus = "failed"
	ShipmentStatusReturned       ShipmentStatus = "returned"
	ShipmentStatusCancelled      ShipmentStatus = "cancelled"
)

type FulfillmentMethod string

const (
	FulfillmentMethodCourier      FulfillmentMethod = "courier"
	FulfillmentMethodSelfDelivery FulfillmentMethod = "self_delivery"
)

type Shipment struct {
	ID uuid.UUID

	OrderID uuid.UUID

	Status            ShipmentStatus
	FulfillmentMethod FulfillmentMethod
	TrackingNumber    *string

	Courier string
	Service string

	Cost          int64
	Weight        int
	OriginID      string
	DestinationID string

	CreatedAt time.Time
}

func (d *Shipment) Validate() error {
	if d.FulfillmentMethod == "" {
		d.FulfillmentMethod = FulfillmentMethodCourier
	}

	if d.FulfillmentMethod != FulfillmentMethodCourier && d.FulfillmentMethod != FulfillmentMethodSelfDelivery {
		return ErrInvalidFulfillmentMethod
	}

	if d.Cost < 0 {
		return ErrInvalidCost
	}

	if d.Weight <= 0 {
		return ErrInvalidWeight
	}

	if d.FulfillmentMethod == FulfillmentMethodCourier {
		if d.Courier == "" {
			return ErrInvalidCourier
		}

		if d.Service == "" {
			return ErrInvalidService
		}
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
	case
		ShipmentStatusCreated,
		ShipmentStatusPacked,
		ShipmentStatusLabelled,
		ShipmentStatusPickedUp,
		ShipmentStatusInTransit,
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusFailed,
		ShipmentStatusReturned,
		ShipmentStatusCancelled:

		return true
	default:
		return false
	}
}
