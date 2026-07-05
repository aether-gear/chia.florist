package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"service-core/internal/modules/audit/domain"
	"service-core/internal/modules/audit/repository"
	"service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type auditLogRepositoryImpl struct{}

func NewAuditLogRepository() repository.AuditLogRepository {
	return &auditLogRepositoryImpl{}
}

func (r *auditLogRepositoryImpl) Save(
	ctx context.Context,
	exec transaction.Executor,
	log domain.AuditLog,
) error {
	metadata, err := json.Marshal(log.Metadata)
	if err != nil {
		return fmt.Errorf("audit log: failed to marshal metadata: %w", err)
	}

	if log.ID == uuid.Nil {
		log.ID = uuid.New()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}

	query := `
		INSERT INTO audit_logs (
			id,
			category,
			action,
			resource,
			resource_id,
			actor_id,
			outcome,
			request_id,
			client_ip,
			metadata,
			created_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err = exec.Exec(ctx, query,
		log.ID,
		log.Category,
		log.Action,
		log.Resource,
		nullableString(log.ResourceID),
		nullableString(log.ActorID),
		log.Outcome,
		nullableString(log.RequestID),
		nullableString(log.ClientIP),
		metadata,
		log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("audit log: failed to save record: %w", err)
	}

	return nil
}

// nullableString converts an empty string to nil so it maps to SQL NULL
// rather than an empty string in nullable text columns.
func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
