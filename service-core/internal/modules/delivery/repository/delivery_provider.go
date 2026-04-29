package repository

import (
	"service-core/internal/modules/delivery/domain"

	"github.com/google/uuid"
)

type CalculateCostInput struct {
	OriginID      string
	DestinationID string
	Weight        int
	Courier       string
}

type CostOption struct {
	ID uuid.UUID

	Cost int64
	ETD  string

	Courier string
	Service string
}

type DeliveryProvider interface {
	CalculateCost(input CalculateCostInput) ([]CostOption, error)
	Track(trackingNumber string) ([]domain.DeliveryEvent, error)
}
