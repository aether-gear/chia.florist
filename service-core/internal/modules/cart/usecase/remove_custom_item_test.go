package usecase

import (
	"context"
	"encoding/json"
	"testing"

	cartDomain "service-core/internal/modules/cart/domain"

	"github.com/google/uuid"
)

func TestRemoveCustomItem_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	customItemID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items: []cartDomain.CartItem{
			{
				ID:           customItemID,
				ItemType:     cartDomain.ItemTypeCustom,
				ProductID:    nil,
				ShopID:       shopID,
				Quantity:     1,
				CustomDesign: json.RawMessage(`{"metadata":{"version":"1.0.0"}}`),
			},
		},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewRemoveCustomItemUsecase(exec, tx, cartR)

	input := RemoveCustomItemInput{
		CustomerID: customerID,
		CartItemID: customItemID,
	}

	err := uc.Execute(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cartR.saved {
		t.Errorf("expected cart to be saved after removal")
	}

	if cartR.cart.Items[0].DeletedAt == nil {
		t.Errorf("expected item to have non-nil DeletedAt timestamp")
	}
}

func TestRemoveCustomItem_NotFound_ItemDoesNotExist(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	customItemID := uuid.New()
	nonExistentItemID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items: []cartDomain.CartItem{
			{
				ID:           customItemID,
				ItemType:     cartDomain.ItemTypeCustom,
				ProductID:    nil,
				ShopID:       shopID,
				Quantity:     1,
				CustomDesign: json.RawMessage(`{"metadata":{"version":"1.0.0"}}`),
			},
		},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewRemoveCustomItemUsecase(exec, tx, cartR)

	input := RemoveCustomItemInput{
		CustomerID: customerID,
		CartItemID: nonExistentItemID,
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Errorf("expected error for non-existent item ID, got nil")
	}
}

func TestRemoveCustomItem_NotFound_CartDoesNotExist(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()

	cartR := &mockSaveCartRepository{cart: nil}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewRemoveCustomItemUsecase(exec, tx, cartR)

	input := RemoveCustomItemInput{
		CustomerID: customerID,
		CartItemID: uuid.New(),
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Errorf("expected error for non-existent cart, got nil")
	}
}
