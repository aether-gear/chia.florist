package repository

import (
	"context"

	appmiddleware "service-core/internal/common/middleware"
	"service-core/internal/modules/authentication/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type AccountRepository interface {
	GetByEmail(
		ctx context.Context,
		exec transaction.Executor,
		email string,
	) (*domain.Account, error)
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Account, error)
	GetByUserID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Account, error)

	ActivateByUserID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) error

	Create(
		ctx context.Context,
		exec transaction.Executor,
		account domain.Account,
	) error
}

type SessionRepository interface {
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.Session, error)

	RevokeByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) error

	UpdateLastActivityByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) error

	Save(
		ctx context.Context,
		exec transaction.Executor,
		session domain.Session,
	) error
}

type VerificationChallengeRepository interface {
	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.VerificationChallenge, error)

	Save(
		ctx context.Context,
		exec transaction.Executor,
		challenge domain.VerificationChallenge,
	) error
}

type RefreshTokenRepository interface {
	Save(
		ctx context.Context,
		exec transaction.Executor,
		challenge domain.RefreshToken,
	) error
}

type TokenService interface {
	Generate(params GenerateTokenParams) (GeneratedToken, error)
	Validate(token string) (*domain.TokenClaims, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type TokenHasher interface {
	Hash(token string) string
	Compare(hash string, token string) bool
}

type Authenticator interface {
	RequireAuth(
		exec transaction.Executor,
		cookie string,
	) appmiddleware.Middleware

	RequireAnyAuth(
		exec transaction.Executor,
		cookies ...string,
	) appmiddleware.Middleware
}
