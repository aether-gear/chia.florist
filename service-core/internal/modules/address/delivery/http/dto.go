package http

import (
	"time"

	"github.com/google/uuid"
)

type saveCustomerAddressRequest struct {
	AddressID    *string `json:"address_id"`
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

type createShopAddressRequest struct {
	Label       string  `json:"label"`
	Phone       *string `json:"phone"`
	IsActive    string  `json:"is_active"`
	ProvinceID  string  `json:"province_id"`
	CityID      string  `json:"city_id"`
	DistrictID  string  `json:"district_id"`
	VillageID   string  `json:"village_id"`
	FullAddress string  `json:"full_address"`
	PostalCode  string  `json:"postal_code"`
}

type customerAddressResponse struct {
	AddressID    uuid.UUID  `json:"address_id"`
	CustomerID   uuid.UUID  `json:"customer_id"`
	ReceiverName string     `json:"receiver_name"`
	Phone        *string    `json:"phone,omitempty"`
	IsDefault    bool       `json:"is_default"`
	ProvinceID   string     `json:"province_id"`
	CityID       string     `json:"city_id"`
	DistrictID   string     `json:"district_id"`
	VillageID    string     `json:"village_id"`
	FullAddress  string     `json:"full_address"`
	PostalCode   string     `json:"postal_code"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

type shopAddressResponse struct {
	ShopID      uuid.UUID  `json:"shop_id"`
	Label       string     `json:"label"`
	Phone       *string    `json:"phone"`
	IsActive    bool       `json:"is_active"`
	ProvinceID  string     `json:"province_id"`
	CityID      string     `json:"city_id"`
	DistrictID  string     `json:"district_id"`
	VillageID   string     `json:"village_id"`
	FullAddress string     `json:"full_address"`
	PostalCode  string     `json:"postal_code"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
