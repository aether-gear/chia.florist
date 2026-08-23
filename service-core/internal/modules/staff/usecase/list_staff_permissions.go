package usecase

import (
	"context"
	"fmt"

	authzDomain "service-core/internal/modules/authorization/domain"
	authzRepo "service-core/internal/modules/authorization/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type ListStaffPermissionsUsecase struct {
	executor transaction.Executor
	permRepo authzRepo.StaffPermissionRepository
}

func NewListStaffPermissionsUsecase(
	executor transaction.Executor,
	permRepo authzRepo.StaffPermissionRepository,
) *ListStaffPermissionsUsecase {
	return &ListStaffPermissionsUsecase{
		executor: executor,
		permRepo: permRepo,
	}
}

func (u *ListStaffPermissionsUsecase) Execute(
	ctx context.Context,
	staffID uuid.UUID,
) ([]authzDomain.StaffPermission, error) {
	perms, err := u.permRepo.ListByStaffID(ctx, u.executor, staffID)
	if err != nil {
		return nil, fmt.Errorf("failed to list staff shop permissions: %w", err)
	}
	if perms == nil {
		perms = []authzDomain.StaffPermission{}
	}

	return perms, nil
}
