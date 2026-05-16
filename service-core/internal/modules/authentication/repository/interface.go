package repository

import (
	"service-core/internal/modules/authentication/domain"

	"github.com/google/uuid"
)

type AccountRepository interface {
	GetByEmail(email string) (*domain.Account, error)
	GetByID(id uuid.UUID) (*domain.Account, error)

	Create(account domain.Account) error
}

type SessionRepository interface {
	Create(session domain.Session) error
}

type VerificationChallengeRepository interface {
	Create(challenge domain.VerificationChallenge) error
}
