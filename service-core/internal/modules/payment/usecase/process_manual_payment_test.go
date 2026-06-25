package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	orderDomain "service-core/internal/modules/order/domain"
	paymentDomain "service-core/internal/modules/payment/domain"

	"github.com/google/uuid"
)

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

	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	items := []orderDomain.OrderItem{
		{
			ProductID: productID,
			ShopID:    shopID,
			Quantity:  3,
		},
	}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	peRepo := &mockPaymentEventRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	usecase := NewProcessManualPaymentUsecase(
		pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo, &mockTransactor{}, &mockExecutor{},
	)

	err := usecase.Execute(ctx, ProcessManualPaymentInput{
		PaymentID: paymentID,
		Action:    "confirm",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payment.Status != paymentDomain.PaymentStatusPaid {
		t.Errorf("expected payment status to be paid, got %v", payment.Status)
	}

	if order.Status != orderDomain.OrderStatusConfirmed {
		t.Errorf("expected order status to be confirmed, got %v", order.Status)
	}

	expectedCommit := fmt.Sprintf("%s-%s-3", productID, shopID)
	if len(iRepo.commits) != 1 || iRepo.commits[0] != expectedCommit {
		t.Errorf("expected commits %v, got %v", []string{expectedCommit}, iRepo.commits)
	}

	if len(paRepo.decremented) != 1 || paRepo.decremented[0] != accountID {
		t.Errorf("expected decremented payment account to be %v, got %v", accountID, paRepo.decremented)
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

	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	items := []orderDomain.OrderItem{
		{
			ProductID: productID,
			ShopID:    shopID,
			Quantity:  3,
		},
	}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	peRepo := &mockPaymentEventRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	usecase := NewProcessManualPaymentUsecase(
		pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo, &mockTransactor{}, &mockExecutor{},
	)

	err := usecase.Execute(ctx, ProcessManualPaymentInput{
		PaymentID: paymentID,
		Action:    "reject",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payment.Status != paymentDomain.PaymentStatusFailed {
		t.Errorf("expected payment status to be failed, got %v", payment.Status)
	}

	if order.Status != orderDomain.OrderStatusCancelled {
		t.Errorf("expected order status to be cancelled, got %v", order.Status)
	}

	expectedRelease := fmt.Sprintf("%s-%s-3", productID, shopID)
	if len(iRepo.releases) != 1 || iRepo.releases[0] != expectedRelease {
		t.Errorf("expected releases %v, got %v", []string{expectedRelease}, iRepo.releases)
	}

	if len(paRepo.decremented) != 1 || paRepo.decremented[0] != accountID {
		t.Errorf("expected decremented payment account to be %v, got %v", accountID, paRepo.decremented)
	}
}

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
	usecase := NewProcessManualPaymentUsecase(
		pRepo, nil, nil, nil, nil, nil, &mockTransactor{}, &mockExecutor{},
	)

	err := usecase.Execute(ctx, ProcessManualPaymentInput{
		PaymentID: paymentID,
		Action:    "confirm",
	})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
