package usecase

import (
	"context"
	"testing"
	"time"

	addressDomain "service-core/internal/modules/address/domain"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type foMockExecutor struct{}

func (m *foMockExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *foMockExecutor) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *foMockExecutor) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return nil
}

type foMockOrderRepo struct {
	capturedParams orderRepo.FindOrderParams
	orders         []orderDomain.Order
	total          int
	err            error
}

func (m *foMockOrderRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*orderDomain.Order, error) {
	return nil, nil
}
func (m *foMockOrderRepo) GetByNumber(_ context.Context, _ transaction.Executor, _ string) (*orderDomain.Order, error) {
	return nil, nil
}
func (m *foMockOrderRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ orderDomain.OrderStatus) error {
	return nil
}
func (m *foMockOrderRepo) UpdateStatusWithSLA(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ orderDomain.OrderStatus, _ *time.Time, _ *time.Time) error {
	return nil
}
func (m *foMockOrderRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.Order) error {
	return nil
}
func (m *foMockOrderRepo) FindOrders(_ context.Context, _ transaction.Executor, params orderRepo.FindOrderParams) ([]orderDomain.Order, int, error) {
	m.capturedParams = params
	return m.orders, m.total, m.err
}
func (m *foMockOrderRepo) SetConfirmedAndExpiry(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ time.Time, _ time.Time) error {
	return nil
}
func (m *foMockOrderRepo) FindExpiredUnfulfilledOrders(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]orderDomain.Order, error) {
	return nil, nil
}

type foMockOrderItemRepo struct{}

func (m *foMockOrderItemRepo) ListByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]orderDomain.OrderItem, error) {
	return nil, nil
}
func (m *foMockOrderItemRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]orderDomain.OrderItem, error) {
	return []orderDomain.OrderItem{}, nil
}
func (m *foMockOrderItemRepo) ListByShipmentID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]orderDomain.OrderItem, error) {
	return nil, nil
}
func (m *foMockOrderItemRepo) SaveBulk(_ context.Context, _ transaction.Executor, _ []orderDomain.OrderItem) error {
	return nil
}
func (m *foMockOrderItemRepo) AssignShipment(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ []uuid.UUID) error {
	return nil
}

type foMockCustomDesignRepo struct{}

func (m *foMockCustomDesignRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.OrderItemCustomDesign) error {
	return nil
}
func (m *foMockCustomDesignRepo) SaveBulk(_ context.Context, _ transaction.Executor, _ []orderDomain.OrderItemCustomDesign) error {
	return nil
}
func (m *foMockCustomDesignRepo) GetByOrderItemID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*orderDomain.OrderItemCustomDesign, error) {
	return nil, nil
}
func (m *foMockCustomDesignRepo) ListByOrderItemIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) (map[uuid.UUID]orderDomain.OrderItemCustomDesign, error) {
	return map[uuid.UUID]orderDomain.OrderItemCustomDesign{}, nil
}
func (m *foMockCustomDesignRepo) ListByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]orderDomain.OrderItemCustomDesign, error) {
	return []orderDomain.OrderItemCustomDesign{}, nil
}

type foMockPaymentRepo struct{}

func (m *foMockPaymentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.Payment, error) {
	return nil, nil
}
func (m *foMockPaymentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.Payment, error) {
	return nil, nil
}
func (m *foMockPaymentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]paymentDomain.Payment, error) {
	return []paymentDomain.Payment{}, nil
}
func (m *foMockPaymentRepo) ListPendingGateway(_ context.Context, _ transaction.Executor, _ time.Time) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *foMockPaymentRepo) ListPastDuePending(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *foMockPaymentRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.Payment) error {
	return nil
}
func (m *foMockPaymentRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ paymentDomain.PaymentStatus) error {
	return nil
}

type foMockPaymentChannelDataRepo struct{}

func (m *foMockPaymentChannelDataRepo) GetByPaymentID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.PaymentChannelData, error) {
	return nil, nil
}
func (m *foMockPaymentChannelDataRepo) ListByPaymentIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) (map[uuid.UUID]*paymentDomain.PaymentChannelData, error) {
	return map[uuid.UUID]*paymentDomain.PaymentChannelData{}, nil
}
func (m *foMockPaymentChannelDataRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentChannelData) error {
	return nil
}

type foMockShipmentRepo struct{}

func (m *foMockShipmentRepo) Create(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}
func (m *foMockShipmentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *foMockShipmentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *foMockShipmentRepo) ListByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]shipmentDomain.Shipment, error) {
	return []shipmentDomain.Shipment{}, nil
}
func (m *foMockShipmentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]shipmentDomain.Shipment, error) {
	return []shipmentDomain.Shipment{}, nil
}
func (m *foMockShipmentRepo) Save(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}
func (m *foMockShipmentRepo) Update(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}

type foMockAddressRepo struct{}

func (m *foMockAddressRepo) CountByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*int, error) {
	zero := 0
	return &zero, nil
}
func (m *foMockAddressRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *foMockAddressRepo) DeleteByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *foMockAddressRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *foMockAddressRepo) GetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *foMockAddressRepo) ListByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *foMockAddressRepo) ListByIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]addressDomain.CustomerAddress, error) {
	return []addressDomain.CustomerAddress{}, nil
}
func (m *foMockAddressRepo) Save(_ context.Context, _ transaction.Executor, _ addressDomain.CustomerAddress) error {
	return nil
}
func (m *foMockAddressRepo) UnsetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}

func TestFindOrdersUsecase_ShopFilter(t *testing.T) {
	mockOrderRepo := &foMockOrderRepo{}
	uc := NewFindOrdersUsecase(
		&foMockExecutor{},
		mockOrderRepo,
		&foMockOrderItemRepo{},
		&foMockCustomDesignRepo{},
		&foMockPaymentRepo{},
		&foMockPaymentChannelDataRepo{},
		&foMockShipmentRepo{},
		&foMockAddressRepo{},
	)

	shopID := uuid.New()
	shopIDs := []uuid.UUID{uuid.New(), uuid.New()}

	input := FindOrdersInput{
		Page:    1,
		Limit:   10,
		ShopID:  &shopID,
		ShopIDs: shopIDs,
	}

	results, total, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 0 || len(results) != 0 {
		t.Errorf("expected empty results, got total=%d len=%d", total, len(results))
	}

	if mockOrderRepo.capturedParams.ShopID == nil || *mockOrderRepo.capturedParams.ShopID != shopID {
		t.Errorf("expected captured ShopID to be %v, got %v", shopID, mockOrderRepo.capturedParams.ShopID)
	}

	if len(mockOrderRepo.capturedParams.ShopIDs) != 2 {
		t.Errorf("expected 2 ShopIDs, got %d", len(mockOrderRepo.capturedParams.ShopIDs))
	}
}

func TestFindOrdersUsecase_StatusFilter_ClearsSingleStatusWhenStatusesProvided(t *testing.T) {
	mockOrderRepo := &foMockOrderRepo{}
	uc := NewFindOrdersUsecase(
		&foMockExecutor{},
		mockOrderRepo,
		&foMockOrderItemRepo{},
		&foMockCustomDesignRepo{},
		&foMockPaymentRepo{},
		&foMockPaymentChannelDataRepo{},
		&foMockShipmentRepo{},
		&foMockAddressRepo{},
	)

	statusStr := "pending"
	statuses := []string{"pending", "confirmed"}
	input := FindOrdersInput{
		Page:     1,
		Limit:    10,
		Status:   &statusStr,
		Statuses: statuses,
	}

	_, _, err := uc.Execute(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mockOrderRepo.capturedParams.Status != nil {
		t.Errorf("expected Status to be nil when Statuses slice is provided, got %v", *mockOrderRepo.capturedParams.Status)
	}
	if len(mockOrderRepo.capturedParams.Statuses) != 2 {
		t.Errorf("expected Statuses slice length 2, got %d", len(mockOrderRepo.capturedParams.Statuses))
	}
}
