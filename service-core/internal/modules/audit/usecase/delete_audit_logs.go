package usecase

import (
	"context"

	"service-core/internal/modules/audit/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type DeleteAuditLogsUsecase struct {
	executor  transaction.Executor
	auditRepo repository.AuditLogRepository
}

func NewDeleteAuditLogsUsecase(
	executor transaction.Executor,
	auditRepo repository.AuditLogRepository,
) *DeleteAuditLogsUsecase {
	return &DeleteAuditLogsUsecase{
		executor:  executor,
		auditRepo: auditRepo,
	}
}

type DeleteAuditLogsInput struct {
	IDs       []uuid.UUID
	DeleteAll bool
}

func (u *DeleteAuditLogsUsecase) Execute(ctx context.Context, input DeleteAuditLogsInput) error {
	if input.DeleteAll {
		return u.auditRepo.DeleteAll(ctx, u.executor)
	}
	return u.auditRepo.DeleteMultiple(ctx, u.executor, input.IDs)
}
