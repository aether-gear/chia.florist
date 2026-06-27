package usecase

import (
	"context"
	"fmt"
	"time"

	apperrors "service-core/internal/common/errors"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzRepo "service-core/internal/modules/authorization/repository"
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
	membershipRepo authzRepo.StaffMembershipRepository
	roleRepo       authzRepo.RoleRepository
}

func NewAddStaffAccountUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo authenRepo.AccountRepository,
	pwHasher authenRepo.PasswordHasher,
	userRepo userRepo.UserRepository,
	membershipRepo authzRepo.StaffMembershipRepository,
	roleRepo authzRepo.RoleRepository,
) *AddStaffAccountUsecase {
	return &AddStaffAccountUsecase{
		executor:       executor,
		transactor:     transactor,
		accountRepo:    accountRepo,
		pwHasher:       pwHasher,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
		roleRepo:       roleRepo,
	}
}

type AddStaffAccountParams struct {
	ActorAccountID uuid.UUID
	ActorStaffID   uuid.UUID
	StaffID        uuid.UUID
	Email          string
	Name           string
	Username       string
	Password       string
	Phone          *string
}

type AddStaffAccountResult struct {
	AccountID uuid.UUID
}

func (u *AddStaffAccountUsecase) Execute(
	ctx context.Context,
	input AddStaffAccountParams,
) error {
	actorMembership, err := u.membershipRepo.
		GetByAccountIDAndStaffID(
			ctx,
			u.executor,
			input.ActorAccountID,
			input.ActorStaffID,
		)
	if err != nil {
		return fmt.Errorf("failed to verify actor membership: %w", err)
	}
	if actorMembership == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	actorRoles, err := u.membershipRepo.
		ListRolesByAccountIDAndStaffID(
			ctx,
			u.executor,
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

	existing, err := u.accountRepo.
		GetByEmail(ctx, u.executor, input.Email)
	if err != nil {
		return fmt.Errorf("failed to check existing account: %w", err)
	}
	if existing != nil {
		return apperrors.NewConflict("an account with this email already exists")
	}

	staffRole, err := u.roleRepo.
		GetByCode(ctx, u.executor, authzDomain.RoleStaff)
	if err != nil {
		return fmt.Errorf("failed to retrieve staff role: %w", err)
	}
	if staffRole == nil {
		return fmt.Errorf("staff role not found in database")
	}

	now := time.Now()
	newUserID := uuid.New()
	newAccountID := uuid.New()

	hash, err := u.pwHasher.Hash(input.Password)
	if err != nil {
		return fmt.Errorf("failed to generate placeholder password: %w", err)
	}

	newUser := userRepo.CreateUserProps{
		ID:        newUserID,
		Name:      input.Name,
		Username:  input.Username,
		Phone:     input.Phone,
		CreatedAt: now,
	}

	newAccount := authenDomain.Account{
		ID:        newAccountID,
		UserID:    newUserID,
		Email:     input.Email,
		Password:  hash,
		Status:    authenDomain.AccountActive,
		Type:      authenDomain.AccountTypeStaff,
		CreatedAt: now,
	}

	newMembership := authzDomain.StaffMembership{
		ID:        uuid.New(),
		StaffID:   input.StaffID,
		AccountID: newAccountID,
		RoleID:    staffRole.ID,
		CreatedBy: input.ActorAccountID,
		CreatedAt: now,
	}

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.userRepo.
				CreateUser(ctx, exec, newUser); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
			if err := u.accountRepo.
				Create(ctx, exec, newAccount); err != nil {
				return fmt.Errorf("failed to create account: %w", err)
			}
			if err := u.membershipRepo.
				Save(ctx, exec, newMembership); err != nil {
				return fmt.Errorf("failed to save membership: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}
