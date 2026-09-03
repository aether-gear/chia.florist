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
) (err error) {
	audit := &applogger.AuditScope{
		Category: "user_action",
		Action:   "delete_account",
		Resource: "customer",
		Metadata: map[string]any{
			"user_id": authCtx.UserID.String(),
		},
	}
	if authCtx.CustomerID != nil {
		audit.SetResourceID(authCtx.CustomerID.String())
		audit.SetMeta("customer_id", authCtx.CustomerID.String())
	}
	defer applogger.TrackAudit(ctx, u.auditLogger, nil, audit, &err)()

	// Must be customer to perform
	// customer account deletion
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("only customer accounts can be deleted")
	}

	userID := authCtx.UserID
	customerID := *authCtx.CustomerID

	return u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		// Delete customer-domain data
		// owned by the customer module.
		if err := u.customerDeletionService.DeleteCustomerRecord(ctx, exec, customerID); err != nil {
			return fmt.Errorf("failed to delete customer record: %w", err)
		}

		// Delete identity and authentication data
		// owned by the authentication module.
		if err := u.userDeletionService.DeleteUserRecord(ctx, exec, userID); err != nil {
			return fmt.Errorf("failed to delete user record: %w", err)
		}

		return nil
	})
}
