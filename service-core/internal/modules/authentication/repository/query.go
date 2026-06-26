package repository

import (
	"time"

	"service-core/internal/modules/authentication/domain"
	authorzDomain "service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

type GenerateTokenParams struct {
	UserID     uuid.UUID
	SessionID  uuid.UUID
	StaffID    *uuid.UUID
	CustomerID *uuid.UUID
	Roles      []authorzDomain.RoleCode

	Type domain.TokenType

	Duration time.Duration
}

type GeneratedToken struct {
	Token     string
	ExpiresAt time.Time
	Type      domain.TokenType
}
