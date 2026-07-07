package repository

import (
	"context"

	"service-core/internal/modules/audit/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

// AuditLogRepository is the persistence interface
// for audit log records.
//
// Defined in the domain so the infra layer
// depends inward, not outward.
type AuditLogRepository interface {
	Save(
		ctx context.Context,
		exec transaction.Executor,
		log domain.AuditLog,
	) error

	Find(
		ctx context.Context,
		exec transaction.Executor,
		params FindAuditLogsParams,
	) ([]domain.AuditLog, int, error)

	GetByID(
		ctx context.Context,
		exec transaction.Executor,
		id uuid.UUID,
	) (*domain.AuditLog, error)

	DeleteMultiple(
		ctx context.Context,
		exec transaction.Executor,
		ids []uuid.UUID,
	) error

	DeleteAll(
		ctx context.Context,
		exec transaction.Executor,
	) error
}
