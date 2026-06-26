package http

import (
	"time"

	"github.com/google/uuid"
)

// customerResponse is used in list endpoints where all items share
// the same flattened structure.
type customerResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Name        string     `json:"name"`
	Username    string     `json:"username"`
	Phone       *string    `json:"phone"`
	AvatarURL   *string    `json:"avatar_url"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// customerProfileResponse is the dedicated detail response for a
// single customer entity, used in the profile endpoint.
type customerProfileResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	Name        string     `json:"name"`
	Username    string     `json:"username"`
	Phone       *string    `json:"phone"`
	AvatarURL   *string    `json:"avatar_url"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
