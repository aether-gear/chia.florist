package domain

import "github.com/google/uuid"

type District struct {
	ID     uuid.UUID
	CityID uuid.UUID
	Name   string
}
