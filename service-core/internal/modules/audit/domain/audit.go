package domain

import (
	"context"
	"time"

	"github.com/google/uuid"

	"service-core/internal/shared/transaction"
)

type (
	// AuditCategory classifies the source of an audit event.
	AuditCategory string

	// AuditOutcome represents the result of an audited action.
	AuditOutcome string
)

const (
	AuditCategoryUserAction  AuditCategory = "user_action"
	AuditCategoryWAFEvent    AuditCategory = "waf_event"
	AuditCategoryThreatIntel AuditCategory = "threat_intel"

	AuditOutcomeSuccess AuditOutcome = "success"
	AuditOutcomeFailure AuditOutcome = "failure"
	AuditOutcomeBlocked AuditOutcome = "blocked"
)

// AuditLog is the core domain entity for an audit record.
//
// Audit logs are append-only, they are never updated after creation.
type AuditLog struct {
	ID uuid.UUID

	Category   AuditCategory
	Action     string
	Resource   string
	ResourceID string

	// ActorID is the authenticated user or service
	// that performed the action.
	//
	// Empty when the action was performed by an
	// unauthenticated request.
	ActorID string

	Outcome AuditOutcome

	// Correlation fields — populated from the
	// request context.
	RequestID string
	ClientIP  string

	// Metadata holds arbitrary structured data
	// specific to the event.
	Metadata map[string]any

	CreatedAt time.Time
}

// AuditLogRepository is the persistence interface
// for audit log records.
//
// Defined in the domain so the infra layer
// depends inward, not outward.
type AuditLogRepository interface {
	Save(
		ctx context.Context,
		exec transaction.Executor,
		log AuditLog,
	) error
}
