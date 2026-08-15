package usecase

import (
	"context"
	"testing"

	apperrors "service-core/internal/common/errors"
	authzDomain "service-core/internal/modules/authorization/domain"
	staffDomain "service-core/internal/modules/staff/domain"

	"github.com/google/uuid"
)

func TestDeleteStaff_Success(t *testing.T) {
	ctx := context.Background()
	actorAccountID := uuid.New()
	actorStaffID := uuid.New()
	targetStaffID := actorStaffID

	membershipRepo := &mockStaffMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   actorStaffID,
			AccountID: actorAccountID,
		},
		roles: []authzDomain.Role{
			{ID: uuid.New(), Code: authzDomain.RoleStaffAdmin, Name: "Staff Admin"},
		},
	}

	staffRepo := &mockStaffRepo{
		staff: &staffDomain.Staff{
			ID:     targetStaffID,
			UserID: uuid.New(),
		},
	}

	auditLogger := &mockAuditLogger{}

	uc := NewDeleteStaffUsecase(
		&mockExecutor{},
		&mockTransactor{},
		staffRepo,
		membershipRepo,
		auditLogger,
	)

	err := uc.Execute(ctx, DeleteStaffInput{
		ActorAccountID: actorAccountID,
		ActorStaffID:   actorStaffID,
		StaffID:        targetStaffID,
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if staffRepo.deleteCalls != 1 {
		t.Errorf("expected 1 staff delete call, got: %d", staffRepo.deleteCalls)
	}
	if membershipRepo.deleteByStaffCalls != 1 {
		t.Errorf("expected 1 membership delete call, got: %d", membershipRepo.deleteByStaffCalls)
	}
	if len(auditLogger.events) != 1 {
		t.Errorf("expected 1 audit event, got: %d", len(auditLogger.events))
	}
}

func TestDeleteStaff_NotFound(t *testing.T) {
	ctx := context.Background()
	actorAccountID := uuid.New()
	actorStaffID := uuid.New()

	membershipRepo := &mockStaffMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   actorStaffID,
			AccountID: actorAccountID,
		},
		roles: []authzDomain.Role{
			{ID: uuid.New(), Code: authzDomain.RoleStaffAdmin, Name: "Staff Admin"},
		},
	}

	staffRepo := &mockStaffRepo{
		staff: nil,
	}

	uc := NewDeleteStaffUsecase(
		&mockExecutor{},
		&mockTransactor{},
		staffRepo,
		membershipRepo,
		&mockAuditLogger{},
	)

	err := uc.Execute(ctx, DeleteStaffInput{
		ActorAccountID: actorAccountID,
		ActorStaffID:   actorStaffID,
		StaffID:        uuid.New(),
	})

	if err == nil {
		t.Fatal("expected error for not found staff, got nil")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok || appErr.StatusCode != 404 {
		t.Fatalf("expected 404 Not Found error, got: %v", err)
	}
}
