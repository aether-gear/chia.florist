package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzRepo "service-core/internal/modules/authorization/repository"
	staffDomain "service-core/internal/modules/staff/domain"
	staffRepo "service-core/internal/modules/staff/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ListStaffAccountsUsecase struct {
	executor       transaction.Executor
	staffRepo      staffRepo.StaffRepository
	membershipRepo authzRepo.StaffMembershipRepository
	auditLogger    applogger.AuditLogger
}

func NewListStaffAccountsUsecase(
	executor transaction.Executor,
	staffRepo staffRepo.StaffRepository,
	membershipRepo authzRepo.StaffMembershipRepository,
	auditLogger applogger.AuditLogger,
) *ListStaffAccountsUsecase {
	return &ListStaffAccountsUsecase{
		executor:       executor,
		staffRepo:      staffRepo,
		membershipRepo: membershipRepo,
		auditLogger:    auditLogger,
	}
}

type ListStaffAccountsParams struct {
	ActorAccountID uuid.UUID
	ActorStaffID   uuid.UUID
	StaffID        uuid.UUID
}

func (u *ListStaffAccountsUsecase) Execute(
	ctx context.Context,
	input ListStaffAccountsParams,
) ([]authzDomain.StaffAccountMember, error) {
	actorMembership, err := u.membershipRepo.GetByAccountIDAndStaffID(ctx, u.executor,
		input.ActorAccountID,
		input.ActorStaffID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify actor membership: %w", err)
	}
	if actorMembership == nil {
		return nil, apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	actorRoles, err := u.membershipRepo.ListRolesByAccountIDAndStaffID(ctx, u.executor,
		input.ActorAccountID,
		input.ActorStaffID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve actor roles: %w", err)
	}

	foundAdmin := false
	for _, role := range actorRoles {
		if role.Code == authzDomain.RoleStaffAdmin {
			foundAdmin = true
			break
		}
	}
	if !foundAdmin {
		return nil, apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	staff, err := u.staffRepo.GetByID(ctx, u.executor, input.StaffID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve staff: %w", err)
	}
	if staff == nil {
		return nil, apperrors.NewNotFound(staffDomain.ErrNotFoundStaff.Error())
	}

	accounts, err := u.membershipRepo.ListAccountsByStaffID(ctx, u.executor, staff.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list staff accounts: %w", err)
	}
	if accounts == nil {
		accounts = []authzDomain.StaffAccountMember{}
	}

	return accounts, nil
}
