package repository

import (
	"context"

	"service-core/internal/modules/shipment/domain"

	"github.com/google/uuid"
)

type ShipmentRepository interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Shipment, error)
	Create(
		ctx context.Context,
		shipment domain.Shipment,
	) error
	Update(
		ctx context.Context,
		shipment domain.Shipment,
	) error
}

type ShipmentEventRepository interface {
	AddTracking(
		ctx context.Context,
		tracking domain.ShipmentEvent,
	) error
	GetTrackingByShipmentID(
		ctx context.Context,
		shipmentID uuid.UUID,
	) ([]domain.ShipmentEvent, error)
}

type ShipmentMethodRepository interface {
	ListActive(ctx context.Context) ([]domain.ShipmentOption, error)
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.ShipmentOption, error)
}

type ShippingCostProvider interface {
	CalculateCost(
		ctx context.Context,
		input CalculateCostInput,
	) ([]CostOption, error)
}

type ShipmentTracker interface {
	Track(
		ctx context.Context,
		trackingNumber string,
	) ([]domain.ShipmentEvent, error)
}
