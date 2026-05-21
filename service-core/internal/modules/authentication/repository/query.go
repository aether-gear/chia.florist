package repository

import (
	"service-core/internal/modules/authentication/domain"
	"time"

	"github.com/google/uuid"
)

type GenerateTokenParams struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	Type      domain.TokenType
	Duration  time.Duration
}

type GeneratedToken struct {
	Token     string
	ExpiresAt time.Time
	Type      domain.TokenType
}
