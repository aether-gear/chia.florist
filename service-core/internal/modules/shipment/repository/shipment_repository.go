package repository

import (
	"service-core/internal/modules/shipment/domain"

	"github.com/google/uuid"
)

type ShipmentRepository interface {
	GetShipmentByID(id uuid.UUID) (*domain.Shipment, error)
	CreateShipment(shipment domain.Shipment) error
	UpdateShipment(shipment domain.Shipment) error

	AddTracking(tracking domain.ShipmentEvent) error
	GetTrackingByShipmentID(shipmentID uuid.UUID) ([]domain.ShipmentEvent, error)

	ListActiveMethod() ([]domain.ShipmentOption, error)
	GetMethodByID(uuid.UUID) (*domain.ShipmentOption, error)
}
