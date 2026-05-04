package http

import (
	"time"

	"github.com/google/uuid"
)

type CreateAddressRequest struct {
	UserID       string  `json:"user_id"`
	ReceiverName string  `json:"receiver_name"`
	Phone        *string `json:"phone"`
	IsDefault    *string `json:"is_default"`
	ProvinceID   string  `json:"province_id"`
	CityID       string  `json:"city_id"`
	DistrictID   string  `json:"district_id"`
	VillageID    string  `json:"village_id"`
	FullAddress  string  `json:"full_address"`
	PostalCode   string  `json:"postal_code"`
}

type AddressResponse struct {
	UserID       uuid.UUID  `json:"user_id"`
	ReceiverName string     `json:"receiver_name"`
	Phone        *string    `json:"phone"`
	IsDefault    bool       `json:"is_default"`
	ProvinceID   string     `json:"province_id"`
	CityID       string     `json:"city_id"`
	DistrictID   string     `json:"district_id"`
	VillageID    string     `json:"village_id"`
	FullAddress  string     `json:"full_address"`
	PostalCode   string     `json:"postal_code"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
