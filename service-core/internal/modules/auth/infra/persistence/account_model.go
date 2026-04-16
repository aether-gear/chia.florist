package persistence

import (
	"time"

	"github.com/google/uuid"
)

type AccountModel struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`

	LastLoginAt *time.Time `db:"last_login_at"`
}
