package domain

import (
	"context"

	"github.com/google/uuid"
)

type authContextKey struct{}

var requestAuthContextKey authContextKey

type AuthContext struct {
	UserID    uuid.UUID
	SessionID uuid.UUID

	MerchantID *uuid.UUID

	TokenType TokenType

	IsAuthenticated bool

	Roles []string
}

func WithAuthContext(ctx context.Context, authCtx *AuthContext) context.Context {
	return context.WithValue(ctx, requestAuthContextKey, authCtx)
}

func GetAuthContext(ctx context.Context) (*AuthContext, bool) {
	authCtx, ok := ctx.Value(requestAuthContextKey).(*AuthContext)
	if !ok {
		return nil, false
	}

	return authCtx, true
}
