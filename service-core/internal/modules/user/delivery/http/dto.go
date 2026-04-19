package http

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Username    string     `json:"username"`
	Phone       *string    `json:"phone"`
	LastLoginAt *time.Time `json:"last_login_at"`
}
