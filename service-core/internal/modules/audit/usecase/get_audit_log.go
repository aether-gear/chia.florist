package usecase

import (
	"context"
	"fmt"

	apperrors "service-core/internal/common/errors"
	"service-core/internal/modules/audit/domain"
	"service-core/internal/modules/audit/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type GetAuditLogUsecase struct {
	executor  transaction.Executor
	auditRepo repository.AuditLogRepository
}

func NewGetAuditLogUsecase(
	executor transaction.Executor,
	auditRepo repository.AuditLogRepository,
) *GetAuditLogUsecase {
	return &GetAuditLogUsecase{
		executor:  executor,
		auditRepo: auditRepo,
	}
}

type GetAuditLogInput struct {
	ID uuid.UUID
}

func (u *GetAuditLogUsecase) Execute(
	ctx context.Context,
	input GetAuditLogInput,
) (*domain.AuditLog, error) {
	log, err := u.auditRepo.
		GetByID(ctx, u.executor, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}
	if log == nil {
		return nil, apperrors.NewNotFound("audit log not found")
	}

	return log, nil
}
