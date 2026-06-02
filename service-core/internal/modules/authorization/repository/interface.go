package repository

import (
	appmiddleware "service-core/internal/common/middleware"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

type ActorService interface {
	Load(userID uuid.UUID) (*domain.Actor, error)
}

type Authorizer interface {
	LoadActor() appmiddleware.Middleware
	RequireAccountType(allowedTypes ...authendomain.AccountType) appmiddleware.Middleware
}
