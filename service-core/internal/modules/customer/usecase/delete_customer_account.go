package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	addressRepo "service-core/internal/modules/address/repository"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	cartRepo "service-core/internal/modules/cart/repository"
	"service-core/internal/modules/customer/repository"
	userRepo "service-core/internal/modules/user/repository"
	transaction "service-core/internal/shared/transaction"
)

type DeleteCustomerAccountUsecase struct {
	transactor   transaction.Transactor
	accountRepo  authenRepo.AccountRepository
	oauthRepo    authenRepo.OAuthConnectionRepository
	sessionRepo  authenRepo.SessionRepository
	userRepo     userRepo.UserRepository
	customerRepo repository.CustomerRepository
	addressRepo  addressRepo.CustomerAddressRepository
	cartRepo     cartRepo.CartRepository
	auditLogger  applogger.AuditLogger
}

func NewDeleteCustomerAccountUsecase(
	transactor transaction.Transactor,
	accountRepo authenRepo.AccountRepository,
	oauthRepo authenRepo.OAuthConnectionRepository,
	sessionRepo authenRepo.SessionRepository,
	userRepo userRepo.UserRepository,
	customerRepo repository.CustomerRepository,
	addressRepo addressRepo.CustomerAddressRepository,
	cartRepo cartRepo.CartRepository,
	auditLogger applogger.AuditLogger,
) *DeleteCustomerAccountUsecase {
	return &DeleteCustomerAccountUsecase{
		transactor:   transactor,
		accountRepo:  accountRepo,
		oauthRepo:    oauthRepo,
		sessionRepo:  sessionRepo,
		userRepo:     userRepo,
		customerRepo: customerRepo,
		addressRepo:  addressRepo,
		cartRepo:     cartRepo,
		auditLogger:  auditLogger,
	}
}

func (u *DeleteCustomerAccountUsecase) Execute(
	ctx context.Context,
	authCtx authenDomain.AuthContext,
) error {
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("only customer accounts can be deleted")
	}

	userID := authCtx.UserID
	customerID := *authCtx.CustomerID

	err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.customerRepo.
			Delete(ctx, exec, customerID); err != nil {
			return fmt.Errorf("failed to soft delete customer: %w", err)
		}

		if err := u.userRepo.
			Delete(ctx, exec, userID); err != nil {
			return fmt.Errorf("failed to soft delete user profile: %w", err)
		}

		if err := u.accountRepo.
			DeleteByUserID(ctx, exec, userID); err != nil {
			return fmt.Errorf("failed to soft delete account: %w", err)
		}

		if err := u.oauthRepo.
			DeleteByUserID(ctx, exec, userID); err != nil {
			return fmt.Errorf("failed to soft delete oauth connections: %w", err)
		}

		if err := u.addressRepo.
			DeleteByCustomerID(ctx, exec, customerID); err != nil {
			return fmt.Errorf("failed to soft delete customer addresses: %w", err)
		}

		if err := u.cartRepo.
			DeleteByCustomerID(ctx, exec, customerID); err != nil {
			return fmt.Errorf("failed to soft delete cart items: %w", err)
		}

		if err := u.sessionRepo.
			RevokeAllByUserID(ctx, exec, userID); err != nil {
			return fmt.Errorf("failed to revoke sessions: %w", err)
		}

		return nil
	})

	if err != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "delete_account",
			Resource: "customer",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{
				"customer_id": customerID.String(),
				"user_id":     userID.String(),
				"error":       err.Error(),
			},
		})
		return err
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category: "user_action",
		Action:   "delete_account",
		Resource: "customer",
		Outcome:  applogger.OutcomeSuccess,
		Metadata: map[string]any{
			"customer_id": customerID.String(),
			"user_id":     userID.String(),
		},
	})

	return nil
}
