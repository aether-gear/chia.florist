package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	authenRepo "service-core/internal/modules/authentication/repository"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzRepo "service-core/internal/modules/authorization/repository"
	staffDomain "service-core/internal/modules/staff/domain"
	staffRepo "service-core/internal/modules/staff/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type DeleteStaffUsecase struct {
	executor            transaction.Executor
	transactor          transaction.Transactor
	staffRepo           staffRepo.StaffRepository
	membershipRepo      authzRepo.StaffMembershipRepository
	userDeletionService authenRepo.UserDeletionService
	auditLogger         applogger.AuditLogger
}

func NewDeleteStaffUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	staffRepo staffRepo.StaffRepository,
	membershipRepo authzRepo.StaffMembershipRepository,
	userDeletionService authenRepo.UserDeletionService,
	auditLogger applogger.AuditLogger,
) *DeleteStaffUsecase {
	return &DeleteStaffUsecase{
		executor:            executor,
		transactor:          transactor,
		staffRepo:           staffRepo,
		membershipRepo:      membershipRepo,
		userDeletionService: userDeletionService,
		auditLogger:         auditLogger,
	}
}

type DeleteStaffInput struct {
	ActorAccountID uuid.UUID
	ActorStaffID   uuid.UUID
	StaffID        uuid.UUID
}

func (u *DeleteStaffUsecase) Execute(
	ctx context.Context,
	input DeleteStaffInput,
) (err error) {
	audit := &applogger.AuditScope{
		Category:   "user_action",
		Action:     "delete_staff",
		Resource:   "staff",
		ResourceID: input.StaffID.String(),
		Metadata:   map[string]any{"staff_id": input.StaffID.String()},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, nil, audit, &err)()

	actorMembership, err := u.membershipRepo.GetByAccountIDAndStaffID(ctx, u.executor,
		input.ActorAccountID,
		input.ActorStaffID,
	)
	if err != nil {
		return fmt.Errorf("failed to verify actor membership: %w", err)
	}
	if actorMembership == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	actorRoles, err := u.membershipRepo.ListRolesByAccountIDAndStaffID(ctx, u.executor,
		input.ActorAccountID,
		input.ActorStaffID,
	)
	if err != nil {
		return fmt.Errorf("failed to retrieve actor roles: %w", err)
	}

	foundAdmin := false
	for _, role := range actorRoles {
		if role.Code == authzDomain.RoleStaffAdmin {
			foundAdmin = true
			break
		}
	}
	if !foundAdmin {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	staff, err := u.staffRepo.GetByID(ctx, u.executor, input.StaffID)
	if err != nil {
		return fmt.Errorf("failed to retrieve staff: %w", err)
	}
	if staff == nil {
		return apperrors.NewNotFound(staffDomain.ErrNotFoundStaff.Error())
	}

	audit.SetMeta("user_id", staff.UserID.String())

	return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.staffRepo.Delete(ctx, exec,
			input.StaffID,
		); err != nil {
			return fmt.Errorf("failed to soft-delete staff: %w", err)
		}

		if err := u.membershipRepo.DeleteByStaffID(ctx, exec,
			input.StaffID,
		); err != nil {
			return fmt.Errorf("failed to delete staff memberships: %w", err)
		}

		if err := u.userDeletionService.DeleteUserRecord(ctx, exec,
			staff.UserID,
		); err != nil {
			return fmt.Errorf("failed to soft-delete user and accounts for staff: %w", err)
		}

		return nil
	})
}
