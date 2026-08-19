package usecase

import (
	"context"
	"testing"

	authzDomain "service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

func TestListStaffPermissionsUsecase(t *testing.T) {
	ctx := context.Background()
	exec := &mockExecutor{}
	permRepo := &mockStaffPermissionRepo{}

	listUC := NewListStaffPermissionsUsecase(exec, permRepo)
	staffID := uuid.New()
	shopID := uuid.New()

	// 1. Initially empty
	items, err := listUC.Execute(ctx, staffID)
	if err != nil {
		t.Fatalf("unexpected error listing permissions: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 permissions, got %d", len(items))
	}

	// 2. Add item to mock repo directly
	permRepo.Save(ctx, exec, authzDomain.StaffPermission{
		StaffID:     staffID,
		ShopID:      shopID,
		Permissions: []string{authzDomain.PermissionShopUpdate},
	})

	// 3. List should return 1 item
	items, err = listUC.Execute(ctx, staffID)
	if err != nil {
		t.Fatalf("unexpected error listing permissions: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 permission, got %d", len(items))
	}
	if items[0].ShopID != shopID {
		t.Fatalf("expected shopID to be %s, got %s", shopID, items[0].ShopID)
	}
}
