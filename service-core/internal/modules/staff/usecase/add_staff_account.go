package usecase

import (
	"context"
	"fmt"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzRepo "service-core/internal/modules/authorization/repository"
	staffDomain "service-core/internal/modules/staff/domain"
	"service-core/internal/modules/staff/repository"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type AddStaffAccountUsecase struct {
	executor       transaction.Executor
	transactor     transaction.Transactor
	accountRepo    authenRepo.AccountRepository
	pwHasher       authenRepo.PasswordHasher
	userRepo       userRepo.UserRepository
	staffRepo      repository.StaffRepository
	membershipRepo authzRepo.StaffMembershipRepository
	roleRepo       authzRepo.RoleRepository
	auditLogger    applogger.AuditLogger
}

func NewAddStaffAccountUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo authenRepo.AccountRepository,
	pwHasher authenRepo.PasswordHasher,
	userRepo userRepo.UserRepository,
	staffRepo repository.StaffRepository,
	membershipRepo authzRepo.StaffMembershipRepository,
	roleRepo authzRepo.RoleRepository,
	auditLogger applogger.AuditLogger,
) *AddStaffAccountUsecase {
	return &AddStaffAccountUsecase{
		executor:       executor,
		transactor:     transactor,
		accountRepo:    accountRepo,
		pwHasher:       pwHasher,
		userRepo:       userRepo,
		staffRepo:      staffRepo,
		membershipRepo: membershipRepo,
		roleRepo:       roleRepo,
		auditLogger:    auditLogger,
	}
}

type AddStaffAccountParams struct {
	ActorAccountID uuid.UUID
	ActorStaffID   uuid.UUID
	StaffID        uuid.UUID
	Email          string
	Password       string
}

func (u *AddStaffAccountUsecase) Execute(
	ctx context.Context,
	input AddStaffAccountParams,
) (err error) {
	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "add_staff_account",
		Resource: "staff_account",
		Metadata: map[string]any{
			"staff_id": input.StaffID.String(),
			"email":    input.Email,
		},
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, nil, audit, &err)()

	if input.Email == "" {
		return apperrors.NewBadRequest("email is required")
	}
	if input.Password == "" {
		return apperrors.NewBadRequest("password is required")
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

	found := false
	for _, role := range actorRoles {
		if role.Code == authzDomain.RoleStaffAdmin {
			found = true
			break
		}
	}
	if !found {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	existingAcc, err := u.accountRepo.GetByEmail(ctx, u.executor, input.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing account: %w", err)
	}
	if existingAcc != nil {
		return apperrors.NewConflict("an account with this email already exists")
	}

	existingStaff, err := u.staffRepo.GetByID(ctx, u.executor, input.StaffID)
	if err != nil {
		return fmt.Errorf("failed to check existing staff: %w", err)
	}
	if existingStaff == nil {
		return apperrors.NewNotFound(staffDomain.ErrNotFoundStaff.Error())
	}

	existingUserAcc, err := u.accountRepo.GetByUserID(ctx, u.executor, existingStaff.UserID)
	if err != nil {
		return fmt.Errorf("failed to check existing staff account: %w", err)
	}
	if existingUserAcc != nil {
		return apperrors.NewConflict("this staff entity already has a bound account (1 account per user limit)")
	}

	staffRole, err := u.roleRepo.GetByCode(ctx, u.executor, authzDomain.RoleStaff)
	if err != nil {
		return fmt.Errorf("failed to retrieve staff role: %w", err)
	}
	if staffRole == nil {
		return fmt.Errorf("staff role not found in database")
	}

	now := appclock.Now()
	newAccountID := uuid.New()
	audit.SetResourceID(newAccountID.String())

	hash, err := u.pwHasher.Hash(input.Password)
	if err != nil {
		return fmt.Errorf("failed to generate placeholder password: %w", err)
	}

	newAccount := authenDomain.Account{
		ID:        newAccountID,
		UserID:    existingStaff.UserID,
		Email:     input.Email,
		Password:  hash,
		Status:    authenDomain.AccountActive,
		Type:      authenDomain.AccountTypeStaff,
		CreatedAt: now,
	}

	newMembership := authzDomain.StaffMembership{
		ID:        uuid.New(),
		StaffID:   existingStaff.ID,
		AccountID: newAccountID,
		RoleID:    staffRole.ID,
		CreatedBy: input.ActorAccountID,
		CreatedAt: now,
	}

	return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.accountRepo.Create(ctx, exec, newAccount); err != nil {
			return fmt.Errorf("failed to create account: %w", err)
		}
		if err := u.membershipRepo.Save(ctx, exec, newMembership); err != nil {
			return fmt.Errorf("failed to save membership: %w", err)
		}
		return nil
	})
}
