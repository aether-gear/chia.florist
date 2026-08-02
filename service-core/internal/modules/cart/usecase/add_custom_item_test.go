package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	cartDomain "service-core/internal/modules/cart/domain"
	cartRepo "service-core/internal/modules/cart/repository"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

type mockTransactor struct {
	transaction.Transactor
}

func (m *mockTransactor) WithinTransaction(ctx context.Context, fn func(transaction.Executor) error) error {
	return fn(&mockExecutor{})
}

type mockSaveCartRepository struct {
	cartRepo.CartRepository
	cart    *cartDomain.Cart
	saved   bool
	saveErr error
	getErr  error
}

func (m *mockSaveCartRepository) GetWithItemsByCustomerID(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) (*cartDomain.Cart, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.cart, nil
}

func (m *mockSaveCartRepository) NewCart(
	ctx context.Context,
	exec transaction.Executor,
	customerID uuid.UUID,
) (*cartDomain.Cart, error) {
	return &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items:      []cartDomain.CartItem{},
	}, nil
}

func (m *mockSaveCartRepository) Save(
	ctx context.Context,
	exec transaction.Executor,
	cart *cartDomain.Cart,
) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.cart = cart
	m.saved = true
	return nil
}

func TestAddCustomItem_Success(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()
	customPayload := json.RawMessage(`{"metadata":{"version":"1.0.0"},"layout":{"physicalSizeId":"medium"}}`)

	cart := &cartDomain.Cart{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items:      []cartDomain.CartItem{},
	}

	cartR := &mockSaveCartRepository{cart: cart}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddCustomItemUsecase(exec, tx, cartR)

	input := AddCustomItemInput{
		CustomerID:     customerID,
		ShopID:         shopID,
		Quantity:       1,
		ProductName:    "Custom Board — Happy Wedding",
		PhysicalSizeID: "medium",
		CustomDesign:   customPayload,
	}

	err := uc.Execute(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cartR.saved {
		t.Errorf("expected cart to be saved")
	}

	if len(cartR.cart.Items) != 1 {
		t.Fatalf("expected 1 item in cart, got %d", len(cartR.cart.Items))
	}

	item := cartR.cart.Items[0]
	if item.ProductVariantType != cartDomain.ProductVariantTypeCustom {
		t.Errorf("expected ProductVariantTypeCustom, got %s", item.ProductVariantType)
	}
	if item.ProductID != nil {
		t.Errorf("expected nil ProductID, got %v", item.ProductID)
	}
	if item.ShopID != shopID {
		t.Errorf("expected shopID %v, got %v", shopID, item.ShopID)
	}
}

func TestAddCustomItem_ValidationError_MissingDesign(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()

	cartR := &mockSaveCartRepository{}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddCustomItemUsecase(exec, tx, cartR)

	input := AddCustomItemInput{
		CustomerID:     customerID,
		ShopID:         shopID,
		Quantity:       1,
		ProductName:    "Custom Board",
		PhysicalSizeID: "medium",
		CustomDesign:   nil,
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Errorf("expected validation error for missing custom_design, got nil")
	}
}

func TestAddCustomItem_ValidationError_MissingPhysicalSize(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()

	cartR := &mockSaveCartRepository{}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddCustomItemUsecase(exec, tx, cartR)

	input := AddCustomItemInput{
		CustomerID:     customerID,
		ShopID:         shopID,
		Quantity:       1,
		ProductName:    "Custom Board",
		PhysicalSizeID: "",
		CustomDesign:   json.RawMessage(`{"metadata":{"version":"1.0.0"}}`),
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Errorf("expected validation error for missing physical_size_id, got nil")
	}
}

func TestAddCustomItem_InvalidShopID(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()

	cartR := &mockSaveCartRepository{}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddCustomItemUsecase(exec, tx, cartR)

	input := AddCustomItemInput{
		CustomerID:     customerID,
		ShopID:         uuid.Nil,
		Quantity:       1,
		ProductName:    "Custom Board",
		PhysicalSizeID: "medium",
		CustomDesign:   json.RawMessage(`{"metadata":{"version":"1.0.0"}}`),
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Errorf("expected error for nil shopID, got nil")
	}
}

func TestAddCustomItem_GetCartFailure(t *testing.T) {
	ctx := context.Background()
	customerID := uuid.New()
	shopID := uuid.New()

	cartR := &mockSaveCartRepository{getErr: errors.New("db error")}
	tx := &mockTransactor{}
	exec := &mockExecutor{}

	uc := NewAddCustomItemUsecase(exec, tx, cartR)

	input := AddCustomItemInput{
		CustomerID:     customerID,
		ShopID:         shopID,
		Quantity:       1,
		ProductName:    "Custom Board",
		PhysicalSizeID: "medium",
		CustomDesign:   json.RawMessage(`{"metadata":{"version":"1.0.0"}}`),
	}

	err := uc.Execute(ctx, input)
	if err == nil {
		t.Errorf("expected error when get cart fails, got nil")
	}
}
