package repository

import (
	"service-core/internal/modules/authentication/domain"
	"time"

	"github.com/google/uuid"
)

type GenerateTokenParams struct {
	UserID    uuid.UUID
	SessionID uuid.UUID

	MerchantID *uuid.UUID

	Type     domain.TokenType
	Duration time.Duration

	Roles []string
}

type GeneratedToken struct {
	Token     string
	ExpiresAt time.Time
	Type      domain.TokenType
}
