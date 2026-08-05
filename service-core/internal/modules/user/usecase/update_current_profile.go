package usecase

import (
	"context"
	"fmt"
	"strings"

	appclock "service-core/internal/common/clock"
	apperrors "service-core/internal/common/errors"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	customerRepo "service-core/internal/modules/customer/repository"
	staffRepo "service-core/internal/modules/staff/repository"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"
)

type UpdateCurrentProfileUsecase struct {
	executor     transaction.Executor
	transactor   transaction.Transactor
	accountRepo  authenRepo.AccountRepository
	customerRepo customerRepo.CustomerRepository
	staffRepo    staffRepo.StaffRepository
	userRepo     userRepo.UserRepository
}

func NewUpdateCurrentProfileUsecase(
	executor transaction.Executor,
	transactor transaction.Transactor,
	accountRepo authenRepo.AccountRepository,
	customerRepo customerRepo.CustomerRepository,
	staffRepo staffRepo.StaffRepository,
	userRepo userRepo.UserRepository,
) *UpdateCurrentProfileUsecase {
	return &UpdateCurrentProfileUsecase{
		executor:     executor,
		transactor:   transactor,
		accountRepo:  accountRepo,
		customerRepo: customerRepo,
		staffRepo:    staffRepo,
		userRepo:     userRepo,
	}
}

type UpdateProfileInput struct {
	Name      *string
	Phone     *string
	AvatarURL *string
}

func (u *UpdateCurrentProfileUsecase) Execute(
	ctx context.Context,
	authCtx authenDomain.AuthContext,
	input UpdateProfileInput,
) (*ProfileResult, error) {
	if input.Name != nil &&
		strings.TrimSpace(*input.Name) == "" {
		return nil, apperrors.NewBadRequest("name is required")
	}

	account, err := u.accountRepo.
		GetByUserID(ctx, u.executor, authCtx.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	if account == nil {
		return nil, apperrors.NewNotFound(string(apperrors.ErrTypeNotFound))
	}

	var result ProfileResult

	err = u.transactor.WithinTransaction(
		ctx,
		func(exec transaction.Executor) error {
			if err := u.userRepo.SaveProfile(
				ctx,
				exec,
				userRepo.SaveProfileProps{
					UserID:    authCtx.UserID,
					Name:      input.Name,
					Phone:     input.Phone,
					AvatarURL: input.AvatarURL,
					UpdatedAt: appclock.Now(),
				},
			); err != nil {
				return fmt.Errorf("failed to save user profile: %w", err)
			}

			switch account.Type {
			case authenDomain.AccountTypeCustomer:
				customerProfile, err := u.customerRepo.
					GetProfileByUserID(ctx, exec, authCtx.UserID)
				if err != nil {
					return fmt.Errorf("failed to retrieve customer profile: %w", err)
				}

				result.Customer = customerProfile

			case authenDomain.AccountTypeStaff:
				staffProfile, err := u.staffRepo.
					GetProfileByUserID(ctx, exec, authCtx.UserID)
				if err != nil {
					return fmt.Errorf("failed to retrieve staff profile: %w", err)
				}

				result.Staff = staffProfile
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
