package http

import (
	"time"

	"github.com/google/uuid"
)

type addMerchantAccountRequest struct {
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
}

type createMerchantRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	LogoUrl     *string `json:"logo_url"`
	BannerUrl   *string `json:"banner_url"`
}

type merchantResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	LogoUrl     *string    `json:"logo_url"`
	BannerUrl   *string    `json:"banner_url"`
	CreatedAt   time.Time  `json:"created_at"`
}
