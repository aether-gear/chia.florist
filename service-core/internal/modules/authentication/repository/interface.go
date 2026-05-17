package repository

import (
	"service-core/internal/modules/authentication/domain"

	"github.com/google/uuid"
)

type AccountRepository interface {
	GetByEmail(email string) (*domain.Account, error)
	GetByID(id uuid.UUID) (*domain.Account, error)

	ActivateByUserID(id uuid.UUID) error

	Create(account domain.Account) error
}

type SessionRepository interface {
	Create(session domain.Session) error
}

type VerificationChallengeRepository interface {
	GetByID(id uuid.UUID) (*domain.VerificationChallenge, error)

	Create(challenge domain.VerificationChallenge) error
	Save(challenge domain.VerificationChallenge) error
}
