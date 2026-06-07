package repository

import (
	"context"
	appmiddleware "service-core/internal/common/middleware"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authorization/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ActorService interface {
	Load(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
	) (*domain.Actor, error)
}

type Authorizer interface {
	LoadActor(
		exec transaction.Executor,
	) appmiddleware.Middleware
	RequireAccountType(allowedTypes ...authendomain.AccountType) appmiddleware.Middleware
}
