package repository

import (
	"context"
	appmiddleware "service-core/internal/common/middleware"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authorization/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type StaffMembershipRepository interface {
	GetByAccountID(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
	) (*domain.StaffMembership, error)

	GetByAccountIDAndStaffID(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
		staffID uuid.UUID,
	) (*domain.StaffMembership, error)

	ListRolesByAccountIDAndStaffID(
		ctx context.Context,
		exec transaction.Executor,
		accountID uuid.UUID,
		staffID uuid.UUID,
	) ([]domain.Role, error)

	Save(
		ctx context.Context,
		exec transaction.Executor,
		membership domain.StaffMembership,
	) error
}

type RoleRepository interface {
	GetByCode(
		ctx context.Context,
		exec transaction.Executor,
		code domain.RoleCode,
	) (*domain.Role, error)
}

type ActorService interface {
	Load(
		ctx context.Context,
		exec transaction.Executor,
		userID uuid.UUID,
		staffID *uuid.UUID,
	) (*domain.Actor, error)
}

type Authorizer interface {
	LoadActor(
		exec transaction.Executor,
	) appmiddleware.Middleware

	RequireAccountType(
		allowedTypes ...authendomain.AccountType,
	) appmiddleware.Middleware
	RequireStaffRole(
		allowedRoles ...domain.RoleCode,
	) appmiddleware.Middleware
}
