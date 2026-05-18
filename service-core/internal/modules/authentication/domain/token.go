package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type TokenClaims struct {
	UserID    string    `json:"user_id"`
	SessionID string    `json:"session_id"`
	Type      TokenType `json:"type"`

	jwt.RegisteredClaims
}
