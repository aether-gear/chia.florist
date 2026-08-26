package usecase

import (
	"context"
	"testing"

	cartDomain "service-core/internal/modules/cart/domain"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	productDomain "service-core/internal/modules/product/domain"

	"github.com/google/uuid"
)

func TestUpdateItemByID_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items:      []cartDomain.CartItem{},
	}
	_ = cart.AddItem(productID, shopID, 2, cartDomain.ItemOptions{Size: "small", Jambul: "none"})
	cartItemID := cart.Items[0].ID

	cartR := &mockSaveCartRepository{cart: cart}
	invR := &mockAddItemInvRepo{
		inv: &inventoryDomain.Inventory{
			ProductID:  productID,
			ShopID:     shopID,
			TotalStock: 20,
		},
	}
	prodR := &mockAddItemProdRepo{
		prod: &productDomain.Product{
			ID:     productID,
			Name:   "Papan Bunga Premium",
			Status: productDomain.ProductStatusActive,
		},
	}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewUpdateItemUsecase(exec, tx, cartR, invR, prodR)

	err := uc.ExecuteByID(ctx, UpdateItemByIDInput{
		CustomerID:  customerID,
		CartItemID:  cartItemID,
		Quantity:    4,
		ItemOptions: cartDomain.ItemOptions{Size: "large", Jambul: "both"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cartR.saved {
		t.Errorf("expected cart to be saved")
	}

	if cart.Items[0].Quantity != 4 {
		t.Errorf("expected quantity 4, got %d", cart.Items[0].Quantity)
	}
	if cart.Items[0].ItemOptions.Size != "large" || cart.Items[0].ItemOptions.Jambul != "both" {
		t.Errorf("expected options large+both, got %+v", cart.Items[0].ItemOptions)
	}
}

func TestUpdateItemByID_InsufficientStock(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items:      []cartDomain.CartItem{},
	}
	_ = cart.AddItem(productID, shopID, 2)
	cartItemID := cart.Items[0].ID

	cartR := &mockSaveCartRepository{cart: cart}
	invR := &mockAddItemInvRepo{
		inv: &inventoryDomain.Inventory{
			ProductID:  productID,
			ShopID:     shopID,
			TotalStock: 3,
		},
	}
	prodR := &mockAddItemProdRepo{
		prod: &productDomain.Product{
			ID:     productID,
			Name:   "Papan Bunga",
			Status: productDomain.ProductStatusActive,
		},
	}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewUpdateItemUsecase(exec, tx, cartR, invR, prodR)

	err := uc.ExecuteByID(ctx, UpdateItemByIDInput{
		CustomerID: customerID,
		CartItemID: cartItemID,
		Quantity:   5, // exceeds available stock 3
	})

	if err == nil {
		t.Fatalf("expected error for insufficient stock, got nil")
	}
}

func TestUpdateItemByID_CartItemNotFound(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items:      []cartDomain.CartItem{},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	invR := &mockAddItemInvRepo{}
	prodR := &mockAddItemProdRepo{}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewUpdateItemUsecase(exec, tx, cartR, invR, prodR)

	err := uc.ExecuteByID(ctx, UpdateItemByIDInput{
		CustomerID: customerID,
		CartItemID: uuid.New(), // random non-existent cart item
		Quantity:   1,
	})

	if err == nil {
		t.Fatalf("expected error for non-existent cart item, got nil")
	}
}
