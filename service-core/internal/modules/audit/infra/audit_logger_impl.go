package infra

import (
	"context"

	appclock "service-core/internal/common/clock"
	applogger "service-core/internal/common/logger"
	"service-core/internal/modules/audit/domain"
	"service-core/internal/modules/audit/repository"
	"service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type dbAuditLogger struct {
	repo   repository.AuditLogRepository
	syslog applogger.Logger
	exec   transaction.Executor
}

// NewDBAuditLogger returns an AuditLogger that persists events to the database
// and mirrors them to the system logger for real-time stdout observability.
func NewDBAuditLogger(
	repo repository.AuditLogRepository,
	syslog applogger.Logger,
	exec transaction.Executor,
) applogger.AuditLogger {
	return &dbAuditLogger{
		repo:   repo,
		syslog: syslog,
		exec:   exec,
	}
}

// Log persists an AuditEvent to the database and mirrors it to the system logger.
// Actor, request ID, and client IP are extracted automatically from ctx —
// callers do not need to populate those fields manually.
func (a *dbAuditLogger) Log(ctx context.Context, event applogger.AuditEvent) {
	record := domain.AuditLog{
		ID:         uuid.New(),
		Category:   domain.AuditCategory(event.Category),
		Action:     event.Action,
		Resource:   event.Resource,
		ResourceID: event.ResourceID,
		ActorID:    applogger.ActorIDFromContext(ctx),
		Outcome:    domain.AuditOutcome(event.Outcome),
		RequestID:  applogger.RequestIDFromContext(ctx),
		ClientIP:   applogger.ClientIPFromContext(ctx),
		Metadata:   event.Metadata,
		CreatedAt:  appclock.Now(),
	}

	if err := a.repo.Save(ctx, a.exec, record); err != nil {
		// A failed audit write must never crash the request.
		// Log the failure via the system logger so it is visible in monitoring.
		a.syslog.Error(ctx, "audit log: failed to persist event",
			applogger.Field{Key: "error", Value: err.Error()},
			applogger.Field{Key: "action", Value: event.Action},
			applogger.Field{Key: "resource", Value: event.Resource},
			applogger.Field{Key: "outcome", Value: event.Outcome},
		)
	}

	// Mirror to stdout as structured JSON regardless of DB outcome.
	// Useful for real-time log aggregation (Datadog, Loki, etc.).
	a.syslog.Info(ctx, "audit_event",
		applogger.Field{Key: "category", Value: event.Category},
		applogger.Field{Key: "action", Value: event.Action},
		applogger.Field{Key: "resource", Value: event.Resource},
		applogger.Field{Key: "resource_id", Value: event.ResourceID},
		applogger.Field{Key: "outcome", Value: event.Outcome},
	)
}
