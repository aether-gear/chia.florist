package usecase

import (
	"context"
	"testing"

	apperrors "service-core/internal/common/errors"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzDomain "service-core/internal/modules/authorization/domain"
	staffDomain "service-core/internal/modules/staff/domain"

	"github.com/google/uuid"
)

func TestRemoveStaffAccount_Success(t *testing.T) {
	ctx := context.Background()
	actorAccountID := uuid.New()
	targetAccountID := uuid.New()
	targetAccountUserID := uuid.New()
	staffOwnerUserID := uuid.New()
	staffID := uuid.New()

	membershipRepo := &mockStaffMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   staffID,
			AccountID: actorAccountID,
		},
		targetMembership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   staffID,
			AccountID: targetAccountID,
		},
		roles: []authzDomain.Role{
			{ID: uuid.New(), Code: authzDomain.RoleStaffAdmin, Name: "Staff Admin"},
		},
	}

	staffRepo := &mockStaffRepo{
		staff: &staffDomain.Staff{
			ID:     staffID,
			UserID: staffOwnerUserID,
		},
	}

	accountRepo := &mockAccountRepo{
		account: &authenDomain.Account{
			ID:     targetAccountID,
			UserID: targetAccountUserID,
		},
	}

	auditLogger := &mockAuditLogger{}

	uc := NewRemoveStaffAccountUsecase(
		&mockExecutor{},
		&mockTransactor{},
		staffRepo,
		membershipRepo,
		accountRepo,
		auditLogger,
	)

	err := uc.Execute(ctx, RemoveStaffAccountInput{
		ActorAccountID: actorAccountID,
		ActorStaffID:   staffID,
		StaffID:        staffID,
		AccountID:      targetAccountID,
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if membershipRepo.deleteByAccCalls != 1 {
		t.Errorf("expected 1 membership delete call, got: %d", membershipRepo.deleteByAccCalls)
	}
	if len(auditLogger.events) != 1 {
		t.Errorf("expected 1 audit event, got: %d", len(auditLogger.events))
	}
}

func TestRemoveStaffAccount_SelfRemovalBlocked(t *testing.T) {
	ctx := context.Background()
	actorAccountID := uuid.New()
	staffID := uuid.New()

	uc := NewRemoveStaffAccountUsecase(
		&mockExecutor{},
		&mockTransactor{},
		&mockStaffRepo{},
		&mockStaffMembershipRepo{},
		&mockAccountRepo{},
		&mockAuditLogger{},
	)

	err := uc.Execute(ctx, RemoveStaffAccountInput{
		ActorAccountID: actorAccountID,
		ActorStaffID:   staffID,
		StaffID:        staffID,
		AccountID:      actorAccountID,
	})

	if err == nil {
		t.Fatal("expected error when removing own account, got nil")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok || appErr.StatusCode != 400 {
		t.Fatalf("expected 400 Bad Request, got: %v", err)
	}
}

func TestRemoveStaffAccount_PrimaryOwnerRemovalBlocked(t *testing.T) {
	ctx := context.Background()
	actorAccountID := uuid.New()
	ownerAccountID := uuid.New()
	ownerUserID := uuid.New()
	staffID := uuid.New()

	membershipRepo := &mockStaffMembershipRepo{
		membership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   staffID,
			AccountID: actorAccountID,
		},
		targetMembership: &authzDomain.StaffMembership{
			ID:        uuid.New(),
			StaffID:   staffID,
			AccountID: ownerAccountID,
		},
		roles: []authzDomain.Role{
			{ID: uuid.New(), Code: authzDomain.RoleStaffAdmin, Name: "Staff Admin"},
		},
	}

	staffRepo := &mockStaffRepo{
		staff: &staffDomain.Staff{
			ID:     staffID,
			UserID: ownerUserID,
		},
	}

	accountRepo := &mockAccountRepo{
		account: &authenDomain.Account{
			ID:     ownerAccountID,
			UserID: ownerUserID, // Matches staff owner
		},
	}

	uc := NewRemoveStaffAccountUsecase(
		&mockExecutor{},
		&mockTransactor{},
		staffRepo,
		membershipRepo,
		accountRepo,
		&mockAuditLogger{},
	)

	err := uc.Execute(ctx, RemoveStaffAccountInput{
		ActorAccountID: actorAccountID,
		ActorStaffID:   staffID,
		StaffID:        staffID,
		AccountID:      ownerAccountID,
	})

	if err == nil {
		t.Fatal("expected error when removing primary staff owner, got nil")
	}

	appErr, ok := err.(*apperrors.AppError)
	if !ok || appErr.StatusCode != 403 {
		t.Fatalf("expected 403 Forbidden, got: %v", err)
	}
}
