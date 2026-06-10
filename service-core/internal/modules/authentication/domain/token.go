package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenClaims struct {
	UserID     uuid.UUID
	SessionID  uuid.UUID
	MerchantID *uuid.UUID

	Type TokenType

	Roles []string

	IssuedAt  time.Time
	ExpiresAt time.Time
}
