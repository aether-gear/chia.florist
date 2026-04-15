package persistence

import (
	"time"

	"github.com/google/uuid"
)

type CartModel struct {
	ID uuid.UUID `db:"id"`

	UserID uuid.UUID `db:"user_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
