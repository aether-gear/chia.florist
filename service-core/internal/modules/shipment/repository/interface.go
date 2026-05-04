package repository

import (
	"service-core/internal/modules/shipment/domain"

	"github.com/google/uuid"
)

type ShipmentRepository interface {
	GetByID(id uuid.UUID) (*domain.Shipment, error)
	Create(shipment domain.Shipment) error
	Update(shipment domain.Shipment) error
}

type ShipmentEventRepository interface {
	AddTracking(tracking domain.ShipmentEvent) error
	GetTrackingByShipmentID(shipmentID uuid.UUID) ([]domain.ShipmentEvent, error)
}

type ShipmentMethodRepository interface {
	ListActive() ([]domain.ShipmentOption, error)
	GetByID(uuid.UUID) (*domain.ShipmentOption, error)
}

type ShipmentProvider interface {
	CalculateCost(input CalculateCostInput) ([]CostOption, error)
	Track(trackingNumber string) ([]domain.ShipmentEvent, error)
}
