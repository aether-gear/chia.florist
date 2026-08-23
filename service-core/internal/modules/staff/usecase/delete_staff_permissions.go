package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	authzRepo "service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type DeleteStaffPermissionUsecase struct {
	transactor transaction.Transactor
	permRepo   authzRepo.StaffPermissionRepository
}

func NewDeleteStaffPermissionUsecase(
	transactor transaction.Transactor,
	permRepo authzRepo.StaffPermissionRepository,
) *DeleteStaffPermissionUsecase {
	return &DeleteStaffPermissionUsecase{
		transactor: transactor,
		permRepo:   permRepo,
	}
}

func (u *DeleteStaffPermissionUsecase) Execute(
	ctx context.Context,
	staffID uuid.UUID,
	shopID uuid.UUID,
) error {
	if staffID == uuid.Nil || shopID == uuid.Nil {
		return apperrors.NewBadRequest("staff_id and shop_id are required")
	}

	if err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		if err := u.permRepo.Delete(ctx, exec, staffID, shopID); err != nil {
			return fmt.Errorf("failed to remove staff permission: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
