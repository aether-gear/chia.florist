package domain_test

import (
	"testing"

	"service-core/internal/modules/cart/domain"

	"github.com/google/uuid"
)

func TestCart_UpdateItemByID_QuantityOnly(t *testing.T) {
	cart := &domain.Cart{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
	}

	productID := uuid.New()
	shopID := uuid.New()
	opt := domain.ItemOptions{Size: "medium", Jambul: "top"}

	err := cart.AddItem(productID, shopID, 2, opt)
	if err != nil {
		t.Fatalf("unexpected error adding item: %v", err)
	}

	targetID := cart.Items[0].ID

	// Update quantity only with same options
	err = cart.UpdateItemByID(targetID, 5, opt)
	if err != nil {
		t.Fatalf("unexpected error updating item: %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(cart.Items))
	}
	if cart.Items[0].Quantity != 5 {
		t.Errorf("expected quantity 5, got %d", cart.Items[0].Quantity)
	}
	if !cart.Items[0].ItemOptions.Equals(opt) {
		t.Errorf("expected options %+v, got %+v", opt, cart.Items[0].ItemOptions)
	}
}

func TestCart_UpdateItemByID_OptionsChange_NoCollision(t *testing.T) {
	cart := &domain.Cart{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
	}

	productID := uuid.New()
	shopID := uuid.New()
	optOld := domain.ItemOptions{Size: "small", Jambul: "none"}
	optNew := domain.ItemOptions{Size: "large", Jambul: "both"}

	_ = cart.AddItem(productID, shopID, 1, optOld)
	targetID := cart.Items[0].ID

	err := cart.UpdateItemByID(targetID, 3, optNew)
	if err != nil {
		t.Fatalf("unexpected error updating item options: %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(cart.Items))
	}
	if cart.Items[0].Quantity != 3 {
		t.Errorf("expected quantity 3, got %d", cart.Items[0].Quantity)
	}
	if !cart.Items[0].ItemOptions.Equals(optNew) {
		t.Errorf("expected options %+v, got %+v", optNew, cart.Items[0].ItemOptions)
	}
}

func TestCart_UpdateItemByID_OptionsChange_WithCollision(t *testing.T) {
	cart := &domain.Cart{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
	}

	productID := uuid.New()
	shopID := uuid.New()
	optA := domain.ItemOptions{Size: "small", Jambul: "none"}
	optB := domain.ItemOptions{Size: "large", Jambul: "both"}

	_ = cart.AddItem(productID, shopID, 2, optA)
	itemAID := cart.Items[0].ID

	_ = cart.AddItem(productID, shopID, 3, optB)
	itemBID := cart.Items[1].ID

	// Change item A's options to optB (collision with item B)
	err := cart.UpdateItemByID(itemAID, 2, optB)
	if err != nil {
		t.Fatalf("unexpected error updating with collision: %v", err)
	}

	// Item A should be soft-deleted, Item B should have qty 2 + 3 = 5
	var itemA, itemB *domain.CartItem
	for i := range cart.Items {
		if cart.Items[i].ID == itemAID {
			itemA = &cart.Items[i]
		}
		if cart.Items[i].ID == itemBID {
			itemB = &cart.Items[i]
		}
	}

	if itemA == nil || itemA.DeletedAt == nil {
		t.Errorf("expected item A to be soft-deleted")
	}
	if itemB == nil || itemB.Quantity != 5 {
		t.Errorf("expected item B quantity to be 5, got %v", itemB)
	}
}

func TestCart_UpdateItemByID_NotFound(t *testing.T) {
	cart := &domain.Cart{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
	}

	err := cart.UpdateItemByID(uuid.New(), 1, domain.ItemOptions{})
	if err != domain.ErrCartItemNotFound {
		t.Errorf("expected ErrCartItemNotFound, got %v", err)
	}
}

func TestCart_UpdateItemByID_InvalidQuantity(t *testing.T) {
	cart := &domain.Cart{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
	}
	pid := uuid.New()
	_ = cart.AddItem(pid, uuid.New(), 1)
	targetID := cart.Items[0].ID

	err := cart.UpdateItemByID(targetID, 0)
	if err != domain.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity for 0, got %v", err)
	}

	err = cart.UpdateItemByID(targetID, -1)
	if err != domain.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity for -1, got %v", err)
	}
}
