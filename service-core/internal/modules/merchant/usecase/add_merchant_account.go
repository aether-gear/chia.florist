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

type AddMerchantAccountUsecase struct {
	accountRepo    authenRepo.AccountRepository
	pwHasher       authenRepo.PasswordHasher
	userRepo       userRepo.UserRepository
	membershipRepo authzRepo.MerchantMembershipRepository
	roleRepo       authzRepo.RoleRepository
	transactor     transaction.Transactor
	executor       transaction.Executor
}

func NewAddMerchantAccountUsecase(
	accountRepo authenRepo.AccountRepository,
	pwHasher authenRepo.PasswordHasher,
	userRepo userRepo.UserRepository,
	membershipRepo authzRepo.MerchantMembershipRepository,
	roleRepo authzRepo.RoleRepository,
	transactor transaction.Transactor,
	executor transaction.Executor,
) *AddMerchantAccountUsecase {
	return &AddMerchantAccountUsecase{
		accountRepo:    accountRepo,
		pwHasher:       pwHasher,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
		roleRepo:       roleRepo,
		transactor:     transactor,
		executor:       executor,
	}
}

type AddMerchantAccountParams struct {
	ActorAccountID  uuid.UUID
	ActorMerchantID uuid.UUID

	MerchantID uuid.UUID
	Email      string
	Name       string
	Username   string
	Password   string
	Phone      *string
}

type AddMerchantAccountResult struct {
	AccountID uuid.UUID
}

func (u *AddMerchantAccountUsecase) Execute(
	ctx context.Context,
	input AddMerchantAccountParams,
) error {
	actorMembership, err := u.membershipRepo.
		GetByAccountIDAndMerchantID(
			ctx,
			u.executor,
			input.ActorAccountID,
			input.ActorMerchantID,
		)
	if err != nil {
		return fmt.Errorf("failed to verify actor membership: %w", err)
	}
	if actorMembership == nil {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	actorRoles, err := u.membershipRepo.
		ListRolesByAccountIDAndMerchantID(
			ctx,
			u.executor,
			input.ActorAccountID,
			input.ActorMerchantID,
		)
	if err != nil {
		return fmt.Errorf("failed to retrieve actor roles: %w", err)
	}

	found := false
	for _, role := range actorRoles {
		if role.Code == authzDomain.RoleMerchantAdmin {
			found = true
			break
		}
	}
	if !found {
		return apperrors.NewForbidden(authzDomain.ErrInsufficientRole.Error())
	}

	existing, err := u.accountRepo.GetByEmail(
		ctx,
		u.executor,
		input.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to check existing account: %w", err)
	}
	if existing != nil {
		return apperrors.NewConflict("an account with this email already exists")
	}

	staffRole, err := u.roleRepo.GetByCode(
		ctx,
		u.executor,
		authzDomain.RoleMerchantStaff,
	)
	if err != nil {
		return fmt.Errorf("failed to retrieve staff role: %w", err)
	}
	if staffRole == nil {
		return fmt.Errorf("merchant_staff role not found in database")
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
		Type:      authenDomain.AccountTypeMerchant,
		CreatedAt: now,
	}

	newMembership := authzDomain.MerchantMembership{
		ID:         uuid.New(),
		MerchantID: input.MerchantID,
		AccountID:  newAccountID,
		RoleID:     staffRole.ID,
		CreatedBy:  input.ActorAccountID,
		CreatedAt:  now,
	}

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.userRepo.CreateUser(ctx, exec, newUser); err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
			if err := u.accountRepo.Create(ctx, exec, newAccount); err != nil {
				return fmt.Errorf("failed to create account: %w", err)
			}
			if err := u.membershipRepo.Save(ctx, exec, newMembership); err != nil {
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
