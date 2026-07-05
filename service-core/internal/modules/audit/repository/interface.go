package repository

import (
	"context"

	"service-core/internal/modules/audit/domain"
	transaction "service-core/internal/shared/transaction"
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
}
