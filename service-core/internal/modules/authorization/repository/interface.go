package repository

import (
	"context"
	appmiddleware "service-core/internal/common/middleware"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authorization/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type MerchantMembershipRepository interface {
	GetByAccountID(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
	) (*domain.MerchantMembership, error)

	GetByAccountIDAndMerchantID(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
		merchantID uuid.UUID,
	) (*domain.MerchantMembership, error)

	ListRolesByAccountIDAndMerchantID(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
		merchantID uuid.UUID,
	) ([]domain.Role, error)

	Save(
		ctx context.Context,
		exec transaction.Executor,
		membership domain.MerchantMembership,
	) error
}

type RoleRepository interface {
	GetByCode(
		ctx context.Context,
		exec transaction.Executor,
		code string,
	) (*domain.Role, error)
}

type ActorService interface {
	Load(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
		merchantID uuid.UUID,
	) (*domain.Actor, error)
}

type Authorizer interface {
	LoadActor(
		exec transaction.Executor,
	) appmiddleware.Middleware
	RequireAccountType(allowedTypes ...authendomain.AccountType) appmiddleware.Middleware
	RequireMerchantRole(allowedRoles ...string) appmiddleware.Middleware
}
