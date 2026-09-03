package applogger

import (
	"context"
	"time"
)

// AuditScope accumulates metadata for a single request scope.
type AuditScope struct {
	Category   string
	Action     string
	Resource   string
	ResourceID string
	Reason     string
	Metadata   map[string]any
}

// SetResourceID updates the target resource ID for the audit scope.
func (s *AuditScope) SetResourceID(id string) {
	s.ResourceID = id
}

// SetReason sets a custom failure or outcome reason on the audit scope.
func (s *AuditScope) SetReason(reason string) {
	s.Reason = reason
}

// SetMeta sets a key-value metadata entry on the audit scope.
func (s *AuditScope) SetMeta(key string, val any) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]any)
	}
	s.Metadata[key] = val
}

// TrackAudit returns a deferred function that automatically records
// both audit events and system operational logs upon function completion.
func TrackAudit(
	ctx context.Context,
	auditLogger AuditLogger,
	sysLogger Logger,
	scope *AuditScope,
	errPtr *error,
) func() {
	start := time.Now()
	return func() {
		if scope == nil {
			return
		}

		err := *errPtr
		outcome := OutcomeSuccess

		meta := make(map[string]any)
		for k, v := range scope.Metadata {
			meta[k] = v
		}

		if err != nil {
			outcome = OutcomeFailure
			meta["error"] = err.Error()
			if scope.Reason != "" {
				meta["reason"] = scope.Reason
			} else {
				meta["reason"] = err.Error()
			}
		}

		// 1. Dispatch Security Audit Event if auditLogger is provided
		if auditLogger != nil {
			auditLogger.Log(ctx, AuditEvent{
				Category:   scope.Category,
				Action:     scope.Action,
				Resource:   scope.Resource,
				ResourceID: scope.ResourceID,
				Outcome:    outcome,
				Metadata:   meta,
			})
		}

		// 2. Dispatch System Operational Log if sysLogger is provided
		if sysLogger != nil {
			duration := time.Since(start)
			if err != nil {
				sysLogger.Error(ctx, scope.Action+" failed",
					Field{Key: "action", Value: scope.Action},
					Field{Key: "resource", Value: scope.Resource},
					Field{Key: "duration_ms", Value: duration.Milliseconds()},
					Field{Key: "error", Value: err.Error()},
				)
			} else {
				sysLogger.Info(ctx, scope.Action+" succeeded",
					Field{Key: "action", Value: scope.Action},
					Field{Key: "resource", Value: scope.Resource},
					Field{Key: "duration_ms", Value: duration.Milliseconds()},
				)
			}
		}
	}
}
