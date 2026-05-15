package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenService interface {
	Generate(userID uuid.UUID) (string, time.Time, error)
	Validate(token string) (uuid.UUID, error)
}
