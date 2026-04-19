package domain

import (
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID uuid.UUID

	UserID uuid.UUID

	ReceiverName string
	Phone        *string

	IsDefault bool

	ProvinceID  string
	CityID      string
	DistrictID  string
	VillageID   string
	FullAddress string
	PostalCode  string

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
