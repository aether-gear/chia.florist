package applogger

import (
	"context"
	"errors"
	"testing"
)

type mockAuditLogger struct {
	lastEvent AuditEvent
	logCount  int
}

func (m *mockAuditLogger) Log(_ context.Context, event AuditEvent) {
	m.lastEvent = event
	m.logCount++
}

func TestTrackAudit_Success(t *testing.T) {
	auditLogger := &mockAuditLogger{}
	ctx := context.Background()
	scope := &AuditScope{
		Category: "user_action",
		Action:   "test_action",
		Resource: "test_resource",
	}

	var testErr error
	func() {
		defer TrackAudit(ctx, auditLogger, nil, scope, &testErr)()
		scope.SetResourceID("res-123")
		scope.SetMeta("key", "val")
	}()

	if auditLogger.logCount != 1 {
		t.Fatalf("expected 1 audit log, got %d", auditLogger.logCount)
	}
	if auditLogger.lastEvent.Outcome != OutcomeSuccess {
		t.Errorf("expected outcome success, got %s", auditLogger.lastEvent.Outcome)
	}
	if auditLogger.lastEvent.ResourceID != "res-123" {
		t.Errorf("expected resource ID res-123, got %s", auditLogger.lastEvent.ResourceID)
	}
	if auditLogger.lastEvent.Metadata["key"] != "val" {
		t.Errorf("expected meta key=val, got %v", auditLogger.lastEvent.Metadata["key"])
	}
}

func TestTrackAudit_Failure(t *testing.T) {
	auditLogger := &mockAuditLogger{}
	ctx := context.Background()
	scope := &AuditScope{
		Category: "user_action",
		Action:   "test_action",
		Resource: "test_resource",
	}

	var testErr error
	func() {
		defer TrackAudit(ctx, auditLogger, nil, scope, &testErr)()
		scope.SetReason("invalid credentials")
		testErr = errors.New("unauthorized")
	}()

	if auditLogger.logCount != 1 {
		t.Fatalf("expected 1 audit log, got %d", auditLogger.logCount)
	}
	if auditLogger.lastEvent.Outcome != OutcomeFailure {
		t.Errorf("expected outcome failure, got %s", auditLogger.lastEvent.Outcome)
	}
	if auditLogger.lastEvent.Metadata["reason"] != "invalid credentials" {
		t.Errorf("expected reason invalid credentials, got %v", auditLogger.lastEvent.Metadata["reason"])
	}
}
