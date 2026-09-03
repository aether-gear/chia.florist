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

type RemoveStaffAccountUsecase struct {
	executor       transaction.Executor
	transactor     transaction.Transactor
	staffRepo      staffRepo.StaffRepository
	membershipRepo authzRepo.StaffMembershipRepository
	accountRepo    authenRepo.AccountRepository
	sessionRepo    authenRepo.SessionRepository
	auditLogger    applogger.AuditLogger
}

func NewRemoveStaffAccountUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	staffRepo staffRepo.StaffRepository,
	membershipRepo authzRepo.StaffMembershipRepository,
	accountRepo authenRepo.AccountRepository,
	sessionRepo authenRepo.SessionRepository,
	auditLogger applogger.AuditLogger,
) *RemoveStaffAccountUsecase {
	return &RemoveStaffAccountUsecase{
		executor:       executor,
		transactor:     transactor,
		staffRepo:      staffRepo,
		membershipRepo: membershipRepo,
		accountRepo:    accountRepo,
		sessionRepo:    sessionRepo,
		auditLogger:    auditLogger,
	}
}

type RemoveStaffAccountInput struct {
	ActorAccountID uuid.UUID
	ActorStaffID   uuid.UUID
	StaffID        uuid.UUID
	AccountID      uuid.UUID
}

func (u *RemoveStaffAccountUsecase) Execute(
	ctx context.Context,
	input RemoveStaffAccountInput,
) (err error) {
	audit := &applogger.AuditScope{
		Category:   "user_action",
		Action:     "remove_staff_account",
		Resource:   "staff_account",
		ResourceID: input.AccountID.String(),
		Metadata: map[string]any{
			"staff_id":   input.StaffID.String(),
			"account_id": input.AccountID.String(),
		},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, nil, audit, &err)()

	if input.ActorAccountID == input.AccountID {
		return apperrors.NewBadRequest("cannot remove own account from staff")
	}

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

	targetMembership, err := u.membershipRepo.GetByAccountIDAndStaffID(ctx, u.executor,
		input.AccountID,
		input.StaffID,
	)
	if err != nil {
		return fmt.Errorf("failed to retrieve target membership: %w", err)
	}
	if targetMembership == nil {
		return apperrors.NewNotFound("staff account membership not found")
	}

	targetAccount, err := u.accountRepo.GetByID(ctx, u.executor, input.AccountID)
	if err != nil {
		return fmt.Errorf("failed to retrieve target account: %w", err)
	}
	if targetAccount == nil {
		return apperrors.NewNotFound("target account not found")
	}

	audit.SetMeta("user_id", targetAccount.UserID.String())

	return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.membershipRepo.DeleteByAccountIDAndStaffID(ctx, exec,
			input.AccountID,
			input.StaffID,
		); err != nil {
			return fmt.Errorf("failed to delete staff membership: %w", err)
		}

		// Simple attempt for deletion.
		// Assuming accont has type staff and at the moment is not bind with oauth.
		if err := u.accountRepo.
			DeleteByUserID(ctx, exec, targetAccount.UserID); err != nil {
			return fmt.Errorf("failed to soft delete account: %w", err)
		}

		if err := u.sessionRepo.
			RevokeAllByUserID(ctx, exec, targetAccount.UserID); err != nil {
			return fmt.Errorf("failed to revoke sessions: %w", err)
		}

		return nil
	})
}
