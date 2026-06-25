package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	paymentgateway "service-core/internal/infra/payment-gateway"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockExecutor struct{}

func (m *mockExecutor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *mockExecutor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *mockExecutor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

type mockTransactor struct{}

func (m *mockTransactor) WithinTransaction(ctx context.Context, fn func(transaction.Executor) error) error {
	return fn(&mockExecutor{})
}

type mockPaymentRepo struct {
	payments map[uuid.UUID]*paymentDomain.Payment
}

func (m *mockPaymentRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*paymentDomain.Payment, error) {
	return m.payments[id], nil
}

func (m *mockPaymentRepo) GetByOrderID(ctx context.Context, exec transaction.Executor, orderID uuid.UUID) (*paymentDomain.Payment, error) {
	for _, p := range m.payments {
		if p.OrderID == orderID {
			return p, nil
		}
	}
	return nil, nil
}

func (m *mockPaymentRepo) UpdateStatus(ctx context.Context, exec transaction.Executor, id uuid.UUID, status paymentDomain.PaymentStatus) error {
	if p, ok := m.payments[id]; ok {
		p.Status = status
		return nil
	}
	return errors.New("not found")
}

func (m *mockPaymentRepo) Save(ctx context.Context, exec transaction.Executor, payment paymentDomain.Payment) error {
	m.payments[payment.ID] = &payment
	return nil
}

type mockPaymentAccountRepo struct {
	decremented []uuid.UUID
}

func (m *mockPaymentAccountRepo) Save(ctx context.Context, exec transaction.Executor, pa paymentDomain.PaymentAccount) error {
	return nil
}
func (m *mockPaymentAccountRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*paymentDomain.PaymentAccount, error) {
	return nil, nil
}
func (m *mockPaymentAccountRepo) RetrieveLeastLoaded(ctx context.Context, exec transaction.Executor, methodID uuid.UUID) (*paymentDomain.PaymentAccount, error) {
	return nil, nil
}
func (m *mockPaymentAccountRepo) IncrementLoad(ctx context.Context, exec transaction.Executor, accountID uuid.UUID) error {
	return nil
}
func (m *mockPaymentAccountRepo) DecrementLoad(ctx context.Context, exec transaction.Executor, accountID uuid.UUID) error {
	m.decremented = append(m.decremented, accountID)
	return nil
}
func (m *mockPaymentAccountRepo) ListByMethodID(ctx context.Context, exec transaction.Executor, methodID uuid.UUID) ([]paymentDomain.PaymentAccount, error) {
	return nil, nil
}
func (m *mockPaymentAccountRepo) ListAll(ctx context.Context, exec transaction.Executor) ([]paymentDomain.PaymentAccount, error) {
	return nil, nil
}

type mockPaymentEventRepo struct {
	events []paymentDomain.PaymentEvent
}

func (m *mockPaymentEventRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*paymentDomain.PaymentEvent, error) {
	return nil, nil
}
func (m *mockPaymentEventRepo) ListByPaymentID(ctx context.Context, exec transaction.Executor, paymentID uuid.UUID) ([]paymentDomain.PaymentEvent, error) {
	return nil, nil
}
func (m *mockPaymentEventRepo) Create(ctx context.Context, exec transaction.Executor, event paymentDomain.PaymentEvent) error {
	m.events = append(m.events, event)
	return nil
}

type mockOrderRepo struct {
	orders map[uuid.UUID]*orderDomain.Order
}

func (m *mockOrderRepo) GetByID(ctx context.Context, exec transaction.Executor, id uuid.UUID) (*orderDomain.Order, error) {
	return m.orders[id], nil
}
func (m *mockOrderRepo) GetByNumber(ctx context.Context, exec transaction.Executor, number string) (*orderDomain.Order, error) {
	for _, o := range m.orders {
		if o.Number == number {
			return o, nil
		}
	}
	return nil, nil
}
func (m *mockOrderRepo) UpdateStatus(ctx context.Context, exec transaction.Executor, id uuid.UUID, status orderDomain.OrderStatus) error {
	if o, ok := m.orders[id]; ok {
		o.Status = status
		return nil
	}
	return errors.New("not found")
}
func (m *mockOrderRepo) Save(ctx context.Context, exec transaction.Executor, order orderDomain.Order) error {
	m.orders[order.ID] = &order
	return nil
}
func (m *mockOrderRepo) FindOrders(ctx context.Context, exec transaction.Executor, params orderRepo.FindOrderParams) ([]orderDomain.Order, int, error) {
	return nil, 0, nil
}

type mockOrderItemRepo struct {
	items map[uuid.UUID][]orderDomain.OrderItem
}

func (m *mockOrderItemRepo) ListByOrderID(ctx context.Context, exec transaction.Executor, orderID uuid.UUID) ([]orderDomain.OrderItem, error) {
	return m.items[orderID], nil
}
func (m *mockOrderItemRepo) ListByOrderIDs(ctx context.Context, exec transaction.Executor, orderIDs []uuid.UUID) ([]orderDomain.OrderItem, error) {
	return nil, nil
}
func (m *mockOrderItemRepo) SaveBulk(ctx context.Context, exec transaction.Executor, items []orderDomain.OrderItem) error {
	return nil
}

type mockInventoryRepo struct {
	commits  []string
	releases []string
}

func (m *mockInventoryRepo) GetByProductIDAndShopID(ctx context.Context, exec transaction.Executor, productID uuid.UUID, shopID uuid.UUID) (*inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *mockInventoryRepo) ListByProductID(ctx context.Context, exec transaction.Executor, productID uuid.UUID) ([]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *mockInventoryRepo) ListByProductIDs(ctx context.Context, exec transaction.Executor, productIDs []uuid.UUID) (map[uuid.UUID][]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *mockInventoryRepo) ListByShopID(ctx context.Context, exec transaction.Executor, shopID uuid.UUID) ([]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *mockInventoryRepo) Create(ctx context.Context, exec transaction.Executor, inventory *inventoryDomain.Inventory) error {
	return nil
}
func (m *mockInventoryRepo) Reserve(ctx context.Context, exec transaction.Executor, productID uuid.UUID, shopID uuid.UUID, qty int) error {
	return nil
}
func (m *mockInventoryRepo) Release(ctx context.Context, exec transaction.Executor, productID uuid.UUID, shopID uuid.UUID, qty int) error {
	m.releases = append(m.releases, fmt.Sprintf("%s-%s-%d", productID, shopID, qty))
	return nil
}
func (m *mockInventoryRepo) Commit(ctx context.Context, exec transaction.Executor, productID uuid.UUID, shopID uuid.UUID, qty int) error {
	m.commits = append(m.commits, fmt.Sprintf("%s-%s-%d", productID, shopID, qty))
	return nil
}

type mockPaymentGateway struct {
	result *paymentgateway.NotificationResult
	err    error
}

func (m *mockPaymentGateway) Charge(ctx context.Context, req paymentgateway.ChargeRequest) (*paymentgateway.ChargeResponse, error) {
	return nil, nil
}
func (m *mockPaymentGateway) ParseNotification(ctx context.Context, payload paymentgateway.NotificationPayload) (*paymentgateway.NotificationResult, error) {
	return m.result, m.err
}
func (m *mockPaymentGateway) CancelTransaction(ctx context.Context, gatewayOrderID string) error {
	return nil
}

func TestProcessPaymentWebhook_Settlement(t *testing.T) {
	ctx := context.Background()

	orderID := uuid.New()
	paymentID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID:        paymentID,
		OrderID:   orderID,
		Status:    paymentDomain.PaymentStatusPending,
		Amount:    100000,
		Provider:  "midtrans",
		CreatedAt: time.Now(),
	}

	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	items := []orderDomain.OrderItem{
		{
			ProductID: productID,
			ShopID:    shopID,
			Quantity:  2,
		},
	}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	peRepo := &mockPaymentEventRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	gateway := &mockPaymentGateway{
		result: &paymentgateway.NotificationResult{
			GatewayOrderID:       orderID.String(),
			GatewayTransactionID: "midtrans-tx-123",
			Status:               paymentgateway.NotificationStatusSettlement,
			RawStatus:            "settlement",
			GrossAmount:          100000,
		},
	}

	usecase := NewProcessPaymentWebhookUsecase(
		pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo, gateway, &mockTransactor{}, &mockExecutor{},
	)

	err := usecase.Execute(ctx, ProcessPaymentWebhookInput{
		Payload: map[string]any{"order_id": orderID.String()},
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

	expectedCommit := fmt.Sprintf("%s-%s-2", productID, shopID)
	if len(iRepo.commits) != 1 || iRepo.commits[0] != expectedCommit {
		t.Errorf("expected commits %v, got %v", []string{expectedCommit}, iRepo.commits)
	}

	if len(peRepo.events) != 1 || peRepo.events[0].EventName != "paid" {
		t.Errorf("expected 1 paid event, got %v", peRepo.events)
	}
}

func TestProcessPaymentWebhook_Expire(t *testing.T) {
	ctx := context.Background()

	orderID := uuid.New()
	paymentID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	payment := &paymentDomain.Payment{
		ID:        paymentID,
		OrderID:   orderID,
		Status:    paymentDomain.PaymentStatusPending,
		Amount:    100000,
		Provider:  "midtrans",
		CreatedAt: time.Now(),
	}

	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	items := []orderDomain.OrderItem{
		{
			ProductID: productID,
			ShopID:    shopID,
			Quantity:  2,
		},
	}

	pRepo := &mockPaymentRepo{payments: map[uuid.UUID]*paymentDomain.Payment{paymentID: payment}}
	paRepo := &mockPaymentAccountRepo{}
	peRepo := &mockPaymentEventRepo{}
	oRepo := &mockOrderRepo{orders: map[uuid.UUID]*orderDomain.Order{orderID: order}}
	oiRepo := &mockOrderItemRepo{items: map[uuid.UUID][]orderDomain.OrderItem{orderID: items}}
	iRepo := &mockInventoryRepo{}

	gateway := &mockPaymentGateway{
		result: &paymentgateway.NotificationResult{
			GatewayOrderID:       orderID.String(),
			GatewayTransactionID: "midtrans-tx-123",
			Status:               paymentgateway.NotificationStatusExpire,
			RawStatus:            "expire",
			GrossAmount:          100000,
		},
	}

	usecase := NewProcessPaymentWebhookUsecase(
		pRepo, paRepo, peRepo, oRepo, oiRepo, iRepo, gateway, &mockTransactor{}, &mockExecutor{},
	)

	err := usecase.Execute(ctx, ProcessPaymentWebhookInput{
		Payload: map[string]any{"order_id": orderID.String()},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if payment.Status != paymentDomain.PaymentStatusExpired {
		t.Errorf("expected payment status to be expired, got %v", payment.Status)
	}

	if order.Status != orderDomain.OrderStatusCancelled {
		t.Errorf("expected order status to be cancelled, got %v", order.Status)
	}

	expectedRelease := fmt.Sprintf("%s-%s-2", productID, shopID)
	if len(iRepo.releases) != 1 || iRepo.releases[0] != expectedRelease {
		t.Errorf("expected releases %v, got %v", []string{expectedRelease}, iRepo.releases)
	}
}
