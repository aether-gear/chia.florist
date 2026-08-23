package usecase

import (
	"context"
	"testing"

	authzDomain "service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

func TestSaveStaffPermissionUsecase(t *testing.T) {
	ctx := context.Background()
	exec := &mockExecutor{}
	transactor := &mockTransactor{}
	permRepo := &mockStaffPermissionRepo{}

	saveUC := NewSaveStaffPermissionUsecase(transactor, permRepo)
	listUC := NewListStaffPermissionsUsecase(exec, permRepo)

	staffID := uuid.New()
	shopID := uuid.New()

	// 1. Validation error when staffID is Nil
	err := saveUC.Execute(ctx, SaveStaffPermissionParams{
		StaffID: uuid.Nil,
		ShopID:  shopID,
	})
	if err == nil {
		t.Fatal("expected error when staff_id is Nil, got nil")
	}

	// 2. Validation error when shopID is Nil
	err = saveUC.Execute(ctx, SaveStaffPermissionParams{
		StaffID: staffID,
		ShopID:  uuid.Nil,
	})
	if err == nil {
		t.Fatal("expected error when shop_id is Nil, got nil")
	}

	// 3. Save valid permissions
	err = saveUC.Execute(ctx, SaveStaffPermissionParams{
		StaffID: staffID,
		ShopID:  shopID,
		Permissions: []string{
			authzDomain.PermissionProductCreate,
			authzDomain.PermissionInventoryManage,
		},
		Rules: map[string]any{
			"can_discount": true,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error saving staff permission: %v", err)
	}

	// Verify item was saved
	items, err := listUC.Execute(ctx, staffID)
	if err != nil {
		t.Fatalf("unexpected error listing permissions: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 saved permission, got %d", len(items))
	}
	if len(items[0].Permissions) != 2 {
		t.Fatalf("expected 2 granted permissions, got %d", len(items[0].Permissions))
	}
}
