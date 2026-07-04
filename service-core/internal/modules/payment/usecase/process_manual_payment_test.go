package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	orderDomain "service-core/internal/modules/order/domain"
	paymentDomain "service-core/internal/modules/payment/domain"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Helpers shared with process_payment_webhook_test.go mocks
// ---------------------------------------------------------------------------

func newManualUsecase(
	pRepo *mockPaymentRepo,
	paRepo *mockPaymentAccountRepo,
	peRepo *mockPaymentEventRepo,
	oRepo *mockOrderRepo,
	oiRepo *mockOrderItemRepo,
	iRepo *mockInventoryRepo,
) *ProcessManualPaymentUsecase {
	return NewProcessManualPaymentUsecase(
		pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo, &mockTransactor{}, &mockExecutor{},
	)
}

// ---------------------------------------------------------------------------
// All terminal states
// ---------------------------------------------------------------------------

func TestProcessManualPayment_Confirm(t *testing.T) {
	ctx := context.Background()

	orderID := uuid.New()
	paymentID := uuid.New()
	accountID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID:               paymentID,
		OrderID:          orderID,
		Status:           paymentDomain.PaymentStatusPending,
		Amount:           100000,
		Provider:         "manual",
		PaymentAccountID: &accountID,
		CreatedAt:        time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	items := []orderDomain.OrderItem{{ProductID: productID, ShopID: shopID, Quantity: 3}}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	peRepo := &mockPaymentEventRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	err := newManualUsecase(pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status != paymentDomain.PaymentStatusPaid {
		t.Errorf("expected payment paid, got %v", payment.Status)
	}
	if order.Status != orderDomain.OrderStatusConfirmed {
		t.Errorf("expected order confirmed, got %v", order.Status)
	}
	expected := fmt.Sprintf("%s-%s-3", productID, shopID)
	if len(iRepo.commits) != 1 || iRepo.commits[0] != expected {
		t.Errorf("expected commits %v, got %v", expected, iRepo.commits)
	}
	if len(paRepo.decremented) != 1 || paRepo.decremented[0] != accountID {
		t.Errorf("expected DecrementLoad for %v, got %v", accountID, paRepo.decremented)
	}
	if len(peRepo.events) != 1 || peRepo.events[0].EventName != "paid" {
		t.Errorf("expected 1 paid event, got %v", peRepo.events)
	}
}

func TestProcessManualPayment_Reject(t *testing.T) {
	ctx := context.Background()

	orderID := uuid.New()
	paymentID := uuid.New()
	accountID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID:               paymentID,
		OrderID:          orderID,
		Status:           paymentDomain.PaymentStatusPending,
		Amount:           100000,
		Provider:         "manual",
		PaymentAccountID: &accountID,
		CreatedAt:        time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	items := []orderDomain.OrderItem{{ProductID: productID, ShopID: shopID, Quantity: 3}}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	peRepo := &mockPaymentEventRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	err := newManualUsecase(pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "reject"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status != paymentDomain.PaymentStatusFailed {
		t.Errorf("expected payment failed, got %v", payment.Status)
	}
	if order.Status != orderDomain.OrderStatusCancelled {
		t.Errorf("expected order cancelled, got %v", order.Status)
	}
	expected := fmt.Sprintf("%s-%s-3", productID, shopID)
	if len(iRepo.releases) != 1 || iRepo.releases[0] != expected {
		t.Errorf("expected releases %v, got %v", expected, iRepo.releases)
	}
	if len(paRepo.decremented) != 1 || paRepo.decremented[0] != accountID {
		t.Errorf("expected DecrementLoad for %v, got %v", accountID, paRepo.decremented)
	}
}

func TestProcessManualPayment_Cancel(t *testing.T) {
	ctx := context.Background()

	orderID := uuid.New()
	paymentID := uuid.New()
	accountID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID:               paymentID,
		OrderID:          orderID,
		Status:           paymentDomain.PaymentStatusPending,
		Amount:           100000,
		Provider:         "manual",
		PaymentAccountID: &accountID,
		CreatedAt:        time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	items := []orderDomain.OrderItem{{ProductID: productID, ShopID: shopID, Quantity: 1}}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	peRepo := &mockPaymentEventRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	// "cancel" is treated identically to "reject"
	err := newManualUsecase(pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "cancel"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status != paymentDomain.PaymentStatusFailed {
		t.Errorf("expected payment failed, got %v", payment.Status)
	}
	if order.Status != orderDomain.OrderStatusCancelled {
		t.Errorf("expected order cancelled, got %v", order.Status)
	}
	if len(iRepo.releases) != 1 {
		t.Errorf("expected 1 inventory release for cancel, got %d", len(iRepo.releases))
	}
}

// ---------------------------------------------------------------------------
// Guard: non-manual provider is rejected
// ---------------------------------------------------------------------------

func TestProcessManualPayment_NonManual(t *testing.T) {
	ctx := context.Background()

	paymentID := uuid.New()
	payment := &paymentDomain.Payment{
		ID:       paymentID,
		OrderID:  uuid.New(),
		Status:   paymentDomain.PaymentStatusPending,
		Provider: "midtrans",
	}
	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}

	err := NewProcessManualPaymentUsecase(
		pRepo, nil, nil, nil, nil, nil, &mockTransactor{}, &mockExecutor{},
	).Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err == nil {
		t.Fatal("expected error for non-manual provider, got nil")
	}
}

// ---------------------------------------------------------------------------
// Idempotency: already-processed payments return conflict
// ---------------------------------------------------------------------------

func TestProcessManualPayment_Idempotent_AlreadyPaid(t *testing.T) {
	ctx := context.Background()
	paymentID := uuid.New()

	payment := &paymentDomain.Payment{
		ID:       paymentID,
		OrderID:  uuid.New(),
		Status:   paymentDomain.PaymentStatusPaid, // already resolved
		Provider: "manual",
		CreatedAt: time.Now(),
	}
	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}

	err := NewProcessManualPaymentUsecase(
		pRepo, nil, nil, nil, nil, nil, &mockTransactor{}, &mockExecutor{},
	).Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err == nil {
		t.Fatal("expected conflict error for already-paid payment, got nil")
	}
}

func TestProcessManualPayment_Idempotent_AlreadyCancelled(t *testing.T) {
	ctx := context.Background()
	paymentID := uuid.New()

	payment := &paymentDomain.Payment{
		ID:       paymentID,
		OrderID:  uuid.New(),
		Status:   paymentDomain.PaymentStatusCancelled,
		Provider: "manual",
		CreatedAt: time.Now(),
	}
	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}

	err := NewProcessManualPaymentUsecase(
		pRepo, nil, nil, nil, nil, nil, &mockTransactor{}, &mockExecutor{},
	).Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "reject"})
	if err == nil {
		t.Fatal("expected conflict error for already-cancelled payment, got nil")
	}
}

// ---------------------------------------------------------------------------
// Aggregate invariant: payment account load always decremented
// ---------------------------------------------------------------------------

func TestProcessManualPayment_InvariantPaymentAccountLoadDecremented(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	paymentID := uuid.New()
	accountID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: orderID, Status: paymentDomain.PaymentStatusPending,
		Amount: 50000, Provider: "manual", PaymentAccountID: &accountID, CreatedAt: time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	items := []orderDomain.OrderItem{{ProductID: productID, ShopID: shopID, Quantity: 1}}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	err := newManualUsecase(pRepo, paRepo, &mockPaymentEventRepo{}, oRepo, oiRepo, iRepo).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paRepo.decremented) != 1 || paRepo.decremented[0] != accountID {
		t.Errorf("expected DecrementLoad for %v, got %v", accountID, paRepo.decremented)
	}
}

func TestProcessManualPayment_InvariantInventoryReleasedOnReject(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	paymentID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: orderID, Status: paymentDomain.PaymentStatusPending,
		Amount: 50000, Provider: "manual", CreatedAt: time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	items := []orderDomain.OrderItem{{ProductID: productID, ShopID: shopID, Quantity: 5}}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	err := newManualUsecase(pRepo, &mockPaymentAccountRepo{}, &mockPaymentEventRepo{}, oRepo, oiRepo, iRepo).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "reject"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := fmt.Sprintf("%s-%s-5", productID, shopID)
	if len(iRepo.releases) != 1 || iRepo.releases[0] != expected {
		t.Errorf("expected inventory release %v, got %v", expected, iRepo.releases)
	}
}

// ---------------------------------------------------------------------------
// Failure recovery
// ---------------------------------------------------------------------------

func TestProcessManualPayment_PaymentNotFound(t *testing.T) {
	ctx := context.Background()
	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{}}

	err := NewProcessManualPaymentUsecase(
		pRepo, nil, nil, nil, nil, nil, &mockTransactor{}, &mockExecutor{},
	).Execute(ctx, ProcessManualPaymentInput{PaymentID: uuid.New(), Action: "confirm"})
	if err == nil {
		t.Fatal("expected error when payment not found, got nil")
	}
}

func TestProcessManualPayment_InvalidAction(t *testing.T) {
	ctx := context.Background()
	paymentID := uuid.New()
	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: uuid.New(), Status: paymentDomain.PaymentStatusPending, Provider: "manual", CreatedAt: time.Now(),
	}
	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}

	err := NewProcessManualPaymentUsecase(
		pRepo, nil, nil, nil, nil, nil, &mockTransactor{}, &mockExecutor{},
	).Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "approve_all_please"})
	if err == nil {
		t.Fatal("expected error for invalid action, got nil")
	}
}

func TestProcessManualPayment_PaymentRepoUpdateFails(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	paymentID := uuid.New()

	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: orderID, Status: paymentDomain.PaymentStatusPending,
		Amount: 50000, Provider: "manual", CreatedAt: time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	pRepo := &mockPaymentRepo{
		payments:  map[uuid.UUID]*paymentDomain.Payment{paymentID: payment},
		updateErr: errors.New("db: update failed"),
	}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{}}

	err := newManualUsecase(pRepo, &mockPaymentAccountRepo{}, &mockPaymentEventRepo{}, oRepo, oiRepo, &mockInventoryRepo{}).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err == nil {
		t.Fatal("expected error when payment repo update fails")
	}
}

func TestProcessManualPayment_OrderRepoUpdateFails(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	paymentID := uuid.New()

	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: orderID, Status: paymentDomain.PaymentStatusPending,
		Amount: 50000, Provider: "manual", CreatedAt: time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	oRepo := &mockOrderRepo{
		orders:    map[uuid.UUID]*orderDomain.Order{orderID: order},
		updateErr: errors.New("db: order update failed"),
	}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{}}

	err := newManualUsecase(pRepo, &mockPaymentAccountRepo{}, &mockPaymentEventRepo{}, oRepo, oiRepo, &mockInventoryRepo{}).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err == nil {
		t.Fatal("expected error when order repo update fails")
	}
}

func TestProcessManualPayment_InventoryCommitFails(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	paymentID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: orderID, Status: paymentDomain.PaymentStatusPending,
		Amount: 50000, Provider: "manual", CreatedAt: time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	items := []orderDomain.OrderItem{{ProductID: productID, ShopID: shopID, Quantity: 1}}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{commitErr: errors.New("inventory: stock underflow")}

	err := newManualUsecase(pRepo, &mockPaymentAccountRepo{}, &mockPaymentEventRepo{}, oRepo, oiRepo, iRepo).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err == nil {
		t.Fatal("expected error when inventory commit fails")
	}
}

func TestProcessManualPayment_InventoryReleaseFails(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	paymentID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: orderID, Status: paymentDomain.PaymentStatusPending,
		Amount: 50000, Provider: "manual", CreatedAt: time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	items := []orderDomain.OrderItem{{ProductID: productID, ShopID: shopID, Quantity: 1}}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{releaseErr: errors.New("inventory: release error")}

	err := newManualUsecase(pRepo, &mockPaymentAccountRepo{}, &mockPaymentEventRepo{}, oRepo, oiRepo, iRepo).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "reject"})
	if err == nil {
		t.Fatal("expected error when inventory release fails")
	}
}

func TestProcessManualPayment_PaymentEventCreateFails(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	paymentID := uuid.New()

	payment := &paymentDomain.Payment{
		ID: paymentID, OrderID: orderID, Status: paymentDomain.PaymentStatusPending,
		Amount: 50000, Provider: "manual", CreatedAt: time.Now(),
	}
	order := &orderDomain.Order{ID: orderID, Status: orderDomain.OrderStatusPending}
	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{}}
	peRepo := &mockPaymentEventRepo{createErr: errors.New("db: event insert failed")}

	err := newManualUsecase(pRepo, &mockPaymentAccountRepo{}, peRepo, oRepo, oiRepo, &mockInventoryRepo{}).
		Execute(ctx, ProcessManualPaymentInput{PaymentID: paymentID, Action: "confirm"})
	if err == nil {
		t.Fatal("expected error when payment event creation fails")
	}
}
