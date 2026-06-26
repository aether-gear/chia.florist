package domain

import (
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer-facing view of a user.
//
// At the moment, customers share the same underlying profile data as users.
// This type exists to support customer-specific queries and use cases while
// allowing customer behavior and attributes to evolve independently in the future.
type Customer struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
