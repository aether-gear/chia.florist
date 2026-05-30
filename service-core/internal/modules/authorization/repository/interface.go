package repository

import (
	"service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

type ActorService interface {
	Load(accountID uuid.UUID) (*domain.Actor, error)
}

type Authorizer interface {
	IsMerchant(actor *domain.Actor) bool
	IsCustomer(actor *domain.Actor) bool
}
