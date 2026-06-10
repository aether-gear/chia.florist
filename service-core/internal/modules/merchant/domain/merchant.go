package domain

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID        uuid.UUID
	AccountID uuid.UUID

	Name        string
	Description *string

	LogoUrl   *string
	BannerUrl *string

	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

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
