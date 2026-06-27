package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	customerDomain "service-core/internal/modules/customer/domain"
	customerRepo "service-core/internal/modules/customer/repository"
	staffDomain "service-core/internal/modules/staff/domain"
	staffRepo "service-core/internal/modules/staff/repository"
	transaction "service-core/internal/shared/transaction"
)

type GetCurrentProfileUsecase struct {
	executor     transaction.Executor
	accountRepo  authenRepo.AccountRepository
	customerRepo customerRepo.CustomerRepository
	staffRepo    staffRepo.StaffRepository
	sessionRepo  authenRepo.SessionRepository
}

func NewGetCurrentProfileUsecase(
	executor transaction.Executor,
	accountRepo authenRepo.AccountRepository,
	customerRepo customerRepo.CustomerRepository,
	staffRepo staffRepo.StaffRepository,
	sessionRepo authenRepo.SessionRepository,
) *GetCurrentProfileUsecase {
	return &GetCurrentProfileUsecase{
		executor:     executor,
		accountRepo:  accountRepo,
		customerRepo: customerRepo,
		staffRepo:    staffRepo,
		sessionRepo:  sessionRepo,
	}
}

type ProfileResult struct {
	Customer *customerDomain.CustomerProfile
	Staff    *staffDomain.StaffProfile
}

func (u *GetCurrentProfileUsecase) Execute(
	ctx context.Context,
	authCtx authenDomain.AuthContext,
) (*ProfileResult, error) {
	account, err := u.accountRepo.
		GetByUserID(ctx, u.executor, authCtx.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve account: %w", err)
	}

	if account == nil {
		return nil, apperrors.NewNotFound(string(apperrors.ErrTypeNotFound))
	}

	session, err := u.sessionRepo.
		GetByID(ctx, u.executor, authCtx.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	var result ProfileResult
	switch account.Type {
	case authenDomain.AccountTypeCustomer:
		customerProfile, err := u.customerRepo.
			GetProfileByUserID(
				ctx,
				u.executor,
				authCtx.UserID,
			)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve customer profile: %w", err)
		}

		result.Customer = customerProfile
		result.Customer.LastLoginAt = session.LastActivityAt
	case authenDomain.AccountTypeStaff:
		staffProfile, err := u.staffRepo.
			GetProfileByUserID(
				ctx,
				u.executor,
				authCtx.UserID,
			)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve staff profile: %w", err)
		}

		result.Staff = staffProfile
		result.Staff.LastLoginAt = session.LastActivityAt
	}

	return &result, nil
}
