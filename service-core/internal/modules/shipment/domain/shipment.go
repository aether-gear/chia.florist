package domain

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type (
	// ShipmentStatus represents the operational status of a
	// shipment throughout the delivery process.
	ShipmentStatus string
	// FulfillmentMethod identifies how an order is fulfilled
	// and delivered to the customer.
	FulfillmentMethod string
)

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

	FulfillmentMethodCourier      FulfillmentMethod = "courier"
	FulfillmentMethodSelfDelivery FulfillmentMethod = "self_delivery"
)

// allowedShipmentTransitions defines the valid state transitions for a
// shipment. The key represents the current shipment status, while the
// associated values represent the statuses that can be transitioned to
// from that state.
var allowedShipmentTransitions = map[ShipmentStatus][]ShipmentStatus{
	ShipmentStatusCreated: {
		ShipmentStatusPacked,
		ShipmentStatusLabelled,
		ShipmentStatusPickedUp,
		ShipmentStatusInTransit,
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusCancelled,
	},
	ShipmentStatusPacked: {
		ShipmentStatusLabelled,
		ShipmentStatusPickedUp,
		ShipmentStatusInTransit,
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusCancelled,
	},
	ShipmentStatusLabelled: {
		ShipmentStatusPickedUp,
		ShipmentStatusInTransit,
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusCancelled,
	},
	ShipmentStatusPickedUp: {
		ShipmentStatusInTransit,
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusFailed,
		ShipmentStatusCancelled,
	},
	ShipmentStatusInTransit: {
		ShipmentStatusOutForDelivery,
		ShipmentStatusDelivered,
		ShipmentStatusFailed,
		ShipmentStatusCancelled,
	},
	ShipmentStatusOutForDelivery: {
		ShipmentStatusDelivered,
		ShipmentStatusFailed,
		ShipmentStatusCancelled,
	},
	ShipmentStatusFailed: {
		ShipmentStatusInTransit,
		ShipmentStatusOutForDelivery,
		ShipmentStatusReturned,
		ShipmentStatusCancelled,
	},
	ShipmentStatusDelivered: {},
	ShipmentStatusReturned:  {},
	ShipmentStatusCancelled: {},
}

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

	Events []ShipmentEvent
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

func (d *Shipment) UpdateStatus(status ShipmentStatus) error {
	if d.Status == status {
		return nil
	}

	if !d.canTransitionTo(status) {
		return fmt.Errorf("invalid shipment status transition: %s → %s", d.Status, status)
	}

	d.Status = status
	return nil
}

func (d *Shipment) canTransitionTo(next ShipmentStatus) bool {
	allowed, exists := allowedShipmentTransitions[d.Status]
	if !exists {
		return false
	}

	return slices.Contains(allowed, next)
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
