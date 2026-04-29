package repository

import (
	"service-core/internal/modules/delivery/domain"

	"github.com/google/uuid"
)

type DeliveryRepository interface {
	GetShipmentByID(id uuid.UUID) (*domain.Delivery, error)
	CreateShipment(shipment domain.Delivery) error
	UpdateShipment(shipment domain.Delivery) error

	AddTracking(tracking domain.DeliveryEvent) error
	GetTrackingByShipmentID(shipmentID uuid.UUID) ([]domain.DeliveryEvent, error)

	ListActiveMethod() ([]domain.DeliveryOption, error)
	GetMethodByID(uuid.UUID) (*domain.DeliveryOption, error)
}
