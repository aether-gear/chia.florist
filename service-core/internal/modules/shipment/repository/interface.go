package repository

import (
	"context"

	"service-core/internal/modules/shipment/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ShipmentRepository interface {
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Shipment, error)
	Create(
		ctx context.Context,
		exec transaction.Executor,
		shipment domain.Shipment,
	) error
	Update(
		ctx context.Context,
		exec transaction.Executor,
		shipment domain.Shipment,
	) error
}

type TrackingRepository interface {
	AddTracking(
		ctx context.Context,
		exec transaction.Executor,
		tracking domain.ShipmentEvent,
	) error
	GetTrackingByShipmentID(
		ctx context.Context,
		exec transaction.Executor,
		shipmentID uuid.UUID,
	) ([]domain.ShipmentEvent, error)
}

type ShipmentMethodRepository interface {
	ListActive(
		ctx context.Context,
		exec transaction.Executor,
	) ([]domain.ShipmentOption, error)
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.ShipmentOption, error)
}

type ShippingCostProvider interface {
	CalculateCost(
		ctx context.Context,
		exec transaction.Executor,
		input CalculateCostInput,
	) ([]CostOption, error)
}

type ShipmentTracker interface {
	Track(
		ctx context.Context,
		exec transaction.Executor,
		trackingNumber string,
	) ([]domain.ShipmentEvent, error)
}
