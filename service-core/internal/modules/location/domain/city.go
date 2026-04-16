package domain

import "github.com/google/uuid"

type City struct {
	ID         uuid.UUID
	ProvinceID uuid.UUID
	Name       string
}
