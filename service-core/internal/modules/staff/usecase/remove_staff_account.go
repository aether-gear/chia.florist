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
	auditLogger    applogger.AuditLogger
}

func NewRemoveStaffAccountUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	staffRepo staffRepo.StaffRepository,
	membershipRepo authzRepo.StaffMembershipRepository,
	accountRepo authenRepo.AccountRepository,
	auditLogger applogger.AuditLogger,
) *RemoveStaffAccountUsecase {
	return &RemoveStaffAccountUsecase{
		executor:       executor,
		transactor:     transactor,
		staffRepo:      staffRepo,
		membershipRepo: membershipRepo,
		accountRepo:    accountRepo,
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
) error {
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
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "remove_staff_account",
			Resource: "staff_account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"staff_id": input.StaffID.String(), "account_id": input.AccountID.String(), "reason": "actor membership not found"},
		})
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
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "remove_staff_account",
			Resource: "staff_account",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"staff_id": input.StaffID.String(), "account_id": input.AccountID.String(), "reason": "actor lacks admin role"},
		})
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

	if u.accountRepo != nil {
		targetAccount, err := u.accountRepo.GetByID(ctx, u.executor, input.AccountID)
		if err == nil &&
			targetAccount != nil &&
			targetAccount.UserID == staff.UserID {
			return apperrors.NewForbidden("cannot remove primary staff owner account")
		}
	}

	err = u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		return u.membershipRepo.DeleteByAccountIDAndStaffID(ctx, exec,
			input.AccountID,
			input.StaffID,
		)
	})
	if err != nil {
		return err
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category:   "user_action",
		Action:     "remove_staff_account",
		Resource:   "staff_account",
		ResourceID: input.AccountID.String(),
		Outcome:    applogger.OutcomeSuccess,
		Metadata:   map[string]any{"staff_id": input.StaffID.String(), "account_id": input.AccountID.String()},
	})

	return nil
}
