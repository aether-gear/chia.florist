package domain

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID        uuid.UUID
	AccountID uuid.UUID

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// Identity
// 	Account ID
// 	Merchant ID
// 	Merchant Name / Display Name
// 	Merchant Slug
// 	Profile Photo
// 	Cover Banner
// 	Description / Bio
// Contact
// 	Email
// 	Phone Number
// 	WhatsApp Number
// 	Customer Service Contact
// 	Address
// 	Country
// 	Province / State
// 	City
// 	District
// 	Postal Code
// 	Full Address
// 	Latitude
// 	Longitude
// Status
// 	Is Active
// 	Is Verified
// 	Is Suspended
// 	Verification Status
// 	Merchant Status (pending, active, rejected)
// Settings
// 	Preferred Language
// 	Preferred Currency
// 	Timezone
// Operational
// 	Opening Hours
// 	Closing Hours
// 	Business Days
// 	Delivery Radius
// 	Pickup Available
// Financial
// 	Bank Account Name
// 	Bank Name
// 	Bank Account Number
// 	E-Wallet Information
// 	Tax Number (optional)
// Audit
// 	Created At
// 	Updated At

// -------

// id
// account_id

// merchant_type
// (name: personal/business)

// display_name
// description

// email
// phone

// address
// city
// province
// postal_code
// latitude
// longitude

// logo_url
// banner_url

// is_active
// is_verified

// bank_name
// bank_account_name
// bank_account_number

// created_at
// updated_at
