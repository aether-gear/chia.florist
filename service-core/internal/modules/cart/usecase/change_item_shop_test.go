package usecase

import (
	"context"
	"encoding/json"
	"testing"

	cartDomain "service-core/internal/modules/cart/domain"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	inventoryRepo "service-core/internal/modules/inventory/repository"
	shopDomain "service-core/internal/modules/shop/domain"
	shopRepo "service-core/internal/modules/shop/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type mockShopRepo struct {
	shopRepo.ShopRepository
	shop *shopDomain.Shop
}

func (m *mockShopRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*shopDomain.Shop, error) {
	if m.shop != nil && m.shop.ID == id {
		return m.shop, nil
	}
	return nil, nil
}

type mockChangeShopInventoryRepo struct {
	inventoryRepo.InventoryRepository
	inventory *inventoryDomain.Inventory
}

func (m *mockChangeShopInventoryRepo) GetByProductIDAndShopID(ctx context.Context, exec transaction.Executor, productID uuid.UUID, shopID uuid.UUID) (*inventoryDomain.Inventory, error) {
	if m.inventory != nil && m.inventory.ProductID == productID && m.inventory.ShopID == shopID {
		return m.inventory, nil
	}
	return nil, nil
}

func TestChangeItemShop_StandardItem_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	oldShopID := uuid.New()
	newShopID := uuid.New()
	cartItemID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items: []cartDomain.CartItem{
			{
				ID:                 cartItemID,
				ProductVariantType: cartDomain.ProductVariantTypeStandard,
				ProductID:          &productID,
				ShopID:             oldShopID,
				Quantity:           2,
			},
		},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	shopR := &mockShopRepo{shop: &shopDomain.Shop{ID: newShopID, IsActive: true, ApprovalStatus: shopDomain.ShopApprovalStatusApproved}}
	invR := &mockChangeShopInventoryRepo{inventory: &inventoryDomain.Inventory{
		ProductID:  productID,
		ShopID:     newShopID,
		TotalStock: 10,
	}}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewChangeItemShopUsecase(exec, tx, cartR, shopR, invR)

	err := uc.Execute(ctx, ChangeItemShopInput{
		CustomerID: customerID,
		CartItemID: cartItemID,
		NewShopID:  newShopID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cart.Items[0].ShopID != newShopID {
		t.Errorf("expected item shop_id to be %v, got %v", newShopID, cart.Items[0].ShopID)
	}
}

func TestChangeItemShop_InsufficientStock_Error(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	oldShopID := uuid.New()
	newShopID := uuid.New()
	cartItemID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items: []cartDomain.CartItem{
			{
				ID:                 cartItemID,
				ProductVariantType: cartDomain.ProductVariantTypeStandard,
				ProductID:          &productID,
				ShopID:             oldShopID,
				Quantity:           5,
			},
		},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	shopR := &mockShopRepo{shop: &shopDomain.Shop{ID: newShopID, IsActive: true, ApprovalStatus: shopDomain.ShopApprovalStatusApproved}}
	// Inventory stock is 2, which is less than item quantity 5
	invR := &mockChangeShopInventoryRepo{inventory: &inventoryDomain.Inventory{
		ProductID:  productID,
		ShopID:     newShopID,
		TotalStock: 2,
	}}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewChangeItemShopUsecase(exec, tx, cartR, shopR, invR)

	err := uc.Execute(ctx, ChangeItemShopInput{
		CustomerID: customerID,
		CartItemID: cartItemID,
		NewShopID:  newShopID,
	})
	if err == nil {
		t.Fatalf("expected error due to insufficient stock, got nil")
	}
}

func TestChangeItemShop_CustomItem_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	oldShopID := uuid.New()
	newShopID := uuid.New()
	cartItemID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items: []cartDomain.CartItem{
			{
				ID:                 cartItemID,
				ProductVariantType: cartDomain.ProductVariantTypeCustom,
				ProductID:          nil,
				ShopID:             oldShopID,
				Quantity:           1,
				CustomDesign:       json.RawMessage(`{}`),
			},
		},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	shopR := &mockShopRepo{shop: &shopDomain.Shop{ID: newShopID, IsActive: true, ApprovalStatus: shopDomain.ShopApprovalStatusApproved}}
	invR := &mockChangeShopInventoryRepo{}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewChangeItemShopUsecase(exec, tx, cartR, shopR, invR)

	err := uc.Execute(ctx, ChangeItemShopInput{
		CustomerID: customerID,
		CartItemID: cartItemID,
		NewShopID:  newShopID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cart.Items[0].ShopID != newShopID {
		t.Errorf("expected custom item shop_id to be %v, got %v", newShopID, cart.Items[0].ShopID)
	}
}

func TestChangeItemShop_MergeWithExistingItem_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	oldShopID := uuid.New()
	newShopID := uuid.New()
	cartItemID1 := uuid.New()
	cartItemID2 := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items: []cartDomain.CartItem{
			{
				ID:                 cartItemID1,
				ProductVariantType: cartDomain.ProductVariantTypeStandard,
				ProductID:          &productID,
				ShopID:             oldShopID,
				Quantity:           2,
			},
			{
				ID:                 cartItemID2,
				ProductVariantType: cartDomain.ProductVariantTypeStandard,
				ProductID:          &productID,
				ShopID:             newShopID,
				Quantity:           3,
			},
		},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	shopR := &mockShopRepo{shop: &shopDomain.Shop{ID: newShopID, IsActive: true, ApprovalStatus: shopDomain.ShopApprovalStatusApproved}}
	invR := &mockChangeShopInventoryRepo{inventory: &inventoryDomain.Inventory{
		ProductID:  productID,
		ShopID:     newShopID,
		TotalStock: 10,
	}}
	exec := &mockExecutor{}
	tx := &mockTransactor{}

	uc := NewChangeItemShopUsecase(exec, tx, cartR, shopR, invR)

	err := uc.Execute(ctx, ChangeItemShopInput{
		CustomerID: customerID,
		CartItemID: cartItemID1,
		NewShopID:  newShopID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First item should be deleted, second item quantity merged (3 + 2 = 5)
	if cart.Items[0].DeletedAt == nil {
		t.Errorf("expected original item to be marked as deleted")
	}
	if cart.Items[1].Quantity != 5 {
		t.Errorf("expected merged quantity 5, got %d", cart.Items[1].Quantity)
	}
}
