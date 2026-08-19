package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzRepo "service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type SaveStaffPermissionParams struct {
	StaffID     uuid.UUID
	ShopID      uuid.UUID
	Permissions []string
	Rules       map[string]any
}

type SaveStaffPermissionUsecase struct {
	transactor transaction.Transactor
	permRepo   authzRepo.StaffPermissionRepository
}

func NewSaveStaffPermissionUsecase(
	transactor transaction.Transactor,
	permRepo authzRepo.StaffPermissionRepository,
) *SaveStaffPermissionUsecase {
	return &SaveStaffPermissionUsecase{
		transactor: transactor,
		permRepo:   permRepo,
	}
}

func (u *SaveStaffPermissionUsecase) Execute(
	ctx context.Context,
	params SaveStaffPermissionParams,
) error {
	if params.StaffID == uuid.Nil {
		return apperrors.NewBadRequest("staff_id is required")
	}
	if params.ShopID == uuid.Nil {
		return apperrors.NewBadRequest("shop_id is required")
	}

	if err := u.transactor.WithinTransaction(ctx, func(exec transaction.Executor) error {
		permission := authzDomain.StaffPermission{
			StaffID:     params.StaffID,
			ShopID:      params.ShopID,
			Permissions: params.Permissions,
			Rules:       params.Rules,
		}
		if permission.Permissions == nil {
			permission.Permissions = []string{}
		}
		if permission.Rules == nil {
			permission.Rules = make(map[string]any)
		}

		if err := u.permRepo.Save(ctx, exec, permission); err != nil {
			return fmt.Errorf("failed to save staff permission: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
