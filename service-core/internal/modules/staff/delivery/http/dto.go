package http

import (
	"time"

	"github.com/google/uuid"
)

type addStaffAccountRequest struct {
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
}

type createStaffRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	LogoUrl     *string `json:"logo_url"`
	BannerUrl   *string `json:"banner_url"`
}

// staffResponse is used in list endpoints.
type staffResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Phone     *string   `json:"phone"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}

// staffProfileResponse is the dedicated detail response used in
// the profile endpoint.
type staffProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Phone     *string   `json:"phone"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
}
