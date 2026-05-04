package repository

import "github.com/google/uuid"

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
