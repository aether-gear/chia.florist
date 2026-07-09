package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	applogger "service-core/internal/common/logger"
	authenDomain "service-core/internal/modules/authentication/domain"
	authenRepo "service-core/internal/modules/authentication/repository"
	customerRepo "service-core/internal/modules/customer/repository"
	transaction "service-core/internal/shared/transaction"
)

type DeleteCustomerAccountUsecase struct {
	transactor              transaction.Transactor
	userDeletionService     authenRepo.UserDeletionService
	customerDeletionService customerRepo.CustomerDeletionService
	auditLogger             applogger.AuditLogger
}

func NewDeleteCustomerAccountUsecase(
	transactor transaction.Transactor,
	userDeletionService authenRepo.UserDeletionService,
	customerDeletionService customerRepo.CustomerDeletionService,
	auditLogger applogger.AuditLogger,
) *DeleteCustomerAccountUsecase {
	return &DeleteCustomerAccountUsecase{
		transactor:              transactor,
		userDeletionService:     userDeletionService,
		customerDeletionService: customerDeletionService,
		auditLogger:             auditLogger,
	}
}

func (u *DeleteCustomerAccountUsecase) Execute(
	ctx context.Context,
	authCtx authenDomain.AuthContext,
) error {
	// Must be customer to perform customer account deletion
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("only customer accounts can be deleted")
	}

	userID := authCtx.UserID
	customerID := *authCtx.CustomerID

	err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		// 1. Delegate customer record & cascading customer data deletion
		if err := u.customerDeletionService.DeleteCustomerRecord(ctx, exec, customerID); err != nil {
			return fmt.Errorf("failed to delete customer record: %w", err)
		}

		// 2. Delegate core user identity record deletion
		if err := u.userDeletionService.DeleteUserRecord(ctx, exec, userID); err != nil {
			return fmt.Errorf("failed to delete user record: %w", err)
		}

		return nil
	})

	if err != nil {
		u.auditLogger.Log(ctx, applogger.AuditEvent{
			Category: "user_action",
			Action:   "delete_account",
			Resource: "customer",
			Outcome:  applogger.OutcomeFailure,
			Metadata: map[string]any{"customer_id": customerID.String(), "user_id": userID.String(), "error": err.Error()},
		})
		return err
	}

	u.auditLogger.Log(ctx, applogger.AuditEvent{
		Category: "user_action",
		Action:   "delete_account",
		Resource: "customer",
		Outcome:  applogger.OutcomeSuccess,
		Metadata: map[string]any{"customer_id": customerID.String(), "user_id": userID.String()},
	})

	return nil
}
