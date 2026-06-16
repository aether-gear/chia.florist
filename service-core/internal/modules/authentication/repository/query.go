package repository

import (
	"time"

	"service-core/internal/modules/authentication/domain"
	authorzDomain "service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

type GenerateTokenParams struct {
	UserID    uuid.UUID
	SessionID uuid.UUID

	MerchantID *uuid.UUID

	Type     domain.TokenType
	Duration time.Duration

	Roles []authorzDomain.RoleCode
}

type GeneratedToken struct {
	Token     string
	ExpiresAt time.Time
	Type      domain.TokenType
}
