package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "service-core/internal/common/errors"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	query "service-core/internal/shared/query"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
)

// --- Mocks for GetPaymentDetail ---

type mockOrderDetailOrderRepo struct {
	order *orderDomain.Order
}

func (m *mockOrderDetailOrderRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*orderDomain.Order, error) {
	if m.order != nil && m.order.ID == id {
		return m.order, nil
	}
	return nil, nil
}
func (m *mockOrderDetailOrderRepo) GetByNumber(_ context.Context, _ transaction.Executor, _ string) (*orderDomain.Order, error) {
	return nil, nil
}
func (m *mockOrderDetailOrderRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ orderDomain.OrderStatus) error {
	return nil
}
func (m *mockOrderDetailOrderRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.Order) error {
	return nil
}
func (m *mockOrderDetailOrderRepo) FindOrders(_ context.Context, _ transaction.Executor, _ orderRepo.FindOrderParams) ([]orderDomain.Order, int, error) {
	return nil, 0, nil
}

type mockOrderDetailInvoiceRepo struct {
	invoice *orderDomain.Invoice
}

func (m *mockOrderDetailInvoiceRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*orderDomain.Invoice, error) {
	return nil, nil
}
func (m *mockOrderDetailInvoiceRepo) GetByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) (*orderDomain.Invoice, error) {
	if m.invoice != nil && m.invoice.OrderID == orderID {
		return m.invoice, nil
	}
	return nil, nil
}
func (m *mockOrderDetailInvoiceRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.Invoice) error {
	return nil
}

type mockOrderDetailPaymentRepo struct {
	payment *paymentDomain.Payment
}

func (m *mockOrderDetailPaymentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.Payment, error) {
	return nil, nil
}
func (m *mockOrderDetailPaymentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) (*paymentDomain.Payment, error) {
	if m.payment != nil && m.payment.OrderID == orderID {
		return m.payment, nil
	}
	return nil, nil
}
func (m *mockOrderDetailPaymentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *mockOrderDetailPaymentRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ paymentDomain.PaymentStatus) error {
	return nil
}
func (m *mockOrderDetailPaymentRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.Payment) error {
	return nil
}
func (m *mockOrderDetailPaymentRepo) ListPendingGateway(_ context.Context, _ transaction.Executor, _ time.Time) ([]paymentDomain.Payment, error) {
	if m.payment != nil && m.payment.Status == paymentDomain.PaymentStatusPending && m.payment.Provider == "gateway" {
		return []paymentDomain.Payment{*m.payment}, nil
	}
	return nil, nil
}
func (m *mockOrderDetailPaymentRepo) ListPastDuePending(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]paymentDomain.Payment, error) {
	return nil, nil
}

type mockOrderDetailPaymentMethodRepo struct {
	method *paymentDomain.PaymentMethod
}

func (m *mockOrderDetailPaymentMethodRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentMethod) error {
	return nil
}
func (m *mockOrderDetailPaymentMethodRepo) FindByName(_ context.Context, _ transaction.Executor, _ string) (*paymentDomain.PaymentMethod, error) {
	return nil, nil
}
func (m *mockOrderDetailPaymentMethodRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*paymentDomain.PaymentMethod, error) {
	if m.method != nil && m.method.ID == id {
		return m.method, nil
	}
	return nil, nil
}
func (m *mockOrderDetailPaymentMethodRepo) ListAll(_ context.Context, _ transaction.Executor, _ query.Sorts) ([]paymentDomain.PaymentMethod, error) {
	return nil, nil
}

type mockOrderDetailPaymentAccountRepo struct{}

type mockOrderDetailPaymentInstructionRepo struct {
	instruction *paymentDomain.PaymentInstruction
}

func (m *mockOrderDetailPaymentInstructionRepo) GetByPaymentMethodID(_ context.Context, _ transaction.Executor, methodID uuid.UUID) (*paymentDomain.PaymentInstruction, error) {
	if m.instruction != nil && m.instruction.PaymentMethodID == methodID {
		return m.instruction, nil
	}
	return nil, nil
}
func (m *mockOrderDetailPaymentInstructionRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentInstruction) error {
	return nil
}

type mockOrderDetailPaymentChannelDataRepo struct {
	channelData *paymentDomain.PaymentChannelData
}

func (m *mockOrderDetailPaymentChannelDataRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentChannelData) error {
	return nil
}
func (m *mockOrderDetailPaymentChannelDataRepo) GetByPaymentID(_ context.Context, _ transaction.Executor, paymentID uuid.UUID) (*paymentDomain.PaymentChannelData, error) {
	if m.channelData != nil && m.channelData.PaymentID == paymentID {
		return m.channelData, nil
	}
	return nil, nil
}
func (m *mockOrderDetailPaymentChannelDataRepo) ListByPaymentIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) (map[uuid.UUID]*paymentDomain.PaymentChannelData, error) {
	return nil, nil
}

// --- Tests ---

func TestGetPaymentDetail_Gateway_Success(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	customerID := uuid.New()
	methodID := uuid.New()
	paymentID := uuid.New()

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
		Number:     "ORD-12345",
		Total:      115000,
	}

	invoice := &orderDomain.Invoice{
		ID:      uuid.New(),
		OrderID: orderID,
		Number:  "INV-12345",
	}

	payment := &paymentDomain.Payment{
		ID:        paymentID,
		OrderID:   orderID,
		MethodID:  methodID,
		Provider:  "midtrans",
		Amount:    115000,
		Status:    paymentDomain.PaymentStatusPending,
		ExpiresAt: pointer(time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)),
	}

	method := &paymentDomain.PaymentMethod{
		ID:   methodID,
		Name: "QRIS",
		Type: paymentDomain.TypeQRCode,
	}

	qrString := "data:image/png;base64,qrstring"
	channelData := &paymentDomain.PaymentChannelData{
		ID:          uuid.New(),
		PaymentID:   paymentID,
		ChannelType: paymentDomain.TypeQRCode,
		DisplayName: "QRIS",
		ActionURL:   &qrString,
	}

	instruction := &paymentDomain.PaymentInstruction{
		ID:              uuid.New(),
		PaymentMethodID: methodID,
		Content:         "Pay invoice {{invoice_number}} amount {{amount}} before {{expired_at}} using {{va_number}}",
	}

	uc := NewGetPaymentDetailUsecase(
		&mockExecutor{},
		&mockOrderDetailOrderRepo{order: order},
		&mockOrderDetailInvoiceRepo{invoice: invoice},
		&mockOrderDetailPaymentRepo{payment: payment},
		&mockOrderDetailPaymentMethodRepo{method: method},
		&mockOrderDetailPaymentInstructionRepo{instruction: instruction},
		&mockOrderDetailPaymentChannelDataRepo{channelData: channelData},
	)

	res, err := uc.Execute(ctx, GetPaymentDetailInput{
		OrderID:    orderID,
		CustomerID: &customerID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res == nil {
		t.Fatal("expected non-nil result")
	}

	if res.Payment.ID != paymentID {
		t.Errorf("expected payment ID %v, got %v", paymentID, res.Payment.ID)
	}

	if res.ChannelData == nil {
		t.Fatal("expected non-nil channel data")
	}

	if res.ChannelData.ChannelType != paymentDomain.TypeQRCode {
		t.Errorf("expected channel type %v, got %v", paymentDomain.TypeQRCode, res.ChannelData.ChannelType)
	}

	if res.Instruction == nil {
		t.Fatal("expected non-nil rendered instruction")
	}

	expectedInstruction := "Pay invoice INV-12345 amount 115000 before 2026-07-12T10:00:00Z using data:image/png;base64,qrstring"
	if *res.Instruction != expectedInstruction {
		t.Errorf("rendered instruction mismatch:\nwant: %q\ngot:  %q", expectedInstruction, *res.Instruction)
	}
}



func TestGetPaymentDetail_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	customerID := uuid.New()

	uc := NewGetPaymentDetailUsecase(
		&mockExecutor{},
		&mockOrderDetailOrderRepo{order: nil},
		&mockOrderDetailInvoiceRepo{},
		&mockOrderDetailPaymentRepo{},
		&mockOrderDetailPaymentMethodRepo{},
		&mockOrderDetailPaymentInstructionRepo{},
		&mockOrderDetailPaymentChannelDataRepo{},
	)

	_, err := uc.Execute(ctx, GetPaymentDetailInput{
		OrderID:    orderID,
		CustomerID: &customerID,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %v", err)
	}

	if appErr.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", appErr.StatusCode)
	}
}

func TestGetPaymentDetail_WrongCustomer(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	customerID := uuid.New()
	wrongCustomerID := uuid.New()

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
	}

	uc := NewGetPaymentDetailUsecase(
		&mockExecutor{},
		&mockOrderDetailOrderRepo{order: order},
		&mockOrderDetailInvoiceRepo{},
		&mockOrderDetailPaymentRepo{},
		&mockOrderDetailPaymentMethodRepo{},
		&mockOrderDetailPaymentInstructionRepo{},
		&mockOrderDetailPaymentChannelDataRepo{},
	)

	_, err := uc.Execute(ctx, GetPaymentDetailInput{
		OrderID:    orderID,
		CustomerID: &wrongCustomerID,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var appErr *apperrors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %v", err)
	}

	if appErr.StatusCode != 404 {
		t.Errorf("expected status code 404 (masked from Forbidden), got %d", appErr.StatusCode)
	}
}

func pointer[T any](v T) *T {
	return &v
}
