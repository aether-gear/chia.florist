package usecase

import (
	"context"
	"testing"

	authzDomain "service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

func TestDeleteStaffPermissionUsecase(t *testing.T) {
	ctx := context.Background()
	exec := &mockExecutor{}
	transactor := &mockTransactor{}
	permRepo := &mockStaffPermissionRepo{}

	saveUC := NewSaveStaffPermissionUsecase(transactor, permRepo)
	deleteUC := NewDeleteStaffPermissionUsecase(transactor, permRepo)
	listUC := NewListStaffPermissionsUsecase(exec, permRepo)

	staffID := uuid.New()
	shopID := uuid.New()

	// 1. Validation error when staffID or shopID is Nil
	err := deleteUC.Execute(ctx, uuid.Nil, shopID)
	if err == nil {
		t.Fatal("expected error when staff_id is Nil, got nil")
	}

	// 2. Save permission first
	err = saveUC.Execute(ctx, SaveStaffPermissionParams{
		StaffID: staffID,
		ShopID:  shopID,
		Permissions: []string{
			authzDomain.PermissionShopUpdate,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error saving permission: %v", err)
	}

	// Verify permission exists
	items, _ := listUC.Execute(ctx, staffID)
	if len(items) != 1 {
		t.Fatalf("expected 1 permission before delete, got %d", len(items))
	}

	// 3. Delete permission
	err = deleteUC.Execute(ctx, staffID, shopID)
	if err != nil {
		t.Fatalf("unexpected error deleting permission: %v", err)
	}

	// Verify list is empty after deletion
	items, _ = listUC.Execute(ctx, staffID)
	if len(items) != 0 {
		t.Fatalf("expected 0 permissions after delete, got %d", len(items))
	}
}
