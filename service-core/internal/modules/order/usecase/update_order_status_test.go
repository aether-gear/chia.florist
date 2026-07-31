package usecase

import (
	"context"
	"testing"

	applogger "service-core/internal/common/logger"
	shipping "service-core/internal/infra/shipping"
	addressDomain "service-core/internal/modules/address/domain"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ===========================================================================
// Mocks for UpdateOrderStatusUsecase tests
// ===========================================================================

type uosMockExecutor struct{}

func (m *uosMockExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *uosMockExecutor) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *uosMockExecutor) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return nil
}

type uosMockTransactor struct{}

func (m *uosMockTransactor) WithinTransaction(ctx context.Context, fn func(transaction.Executor) error) error {
	return fn(&uosMockExecutor{})
}

type uosMockOrderRepo struct {
	order        *orderDomain.Order
	err          error
	updatedState struct {
		id     uuid.UUID
		status orderDomain.OrderStatus
	}
}

func (m *uosMockOrderRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*orderDomain.Order, error) {
	if m.order != nil && m.order.ID == id {
		return m.order, nil
	}
	return nil, m.err
}

func (m *uosMockOrderRepo) GetByNumber(_ context.Context, _ transaction.Executor, _ string) (*orderDomain.Order, error) {
	return nil, nil
}

func (m *uosMockOrderRepo) UpdateStatus(_ context.Context, _ transaction.Executor, id uuid.UUID, status orderDomain.OrderStatus) error {
	m.updatedState.id = id
	m.updatedState.status = status
	return nil
}

func (m *uosMockOrderRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.Order) error {
	return nil
}

func (m *uosMockOrderRepo) FindOrders(_ context.Context, _ transaction.Executor, _ orderRepo.FindOrderParams) ([]orderDomain.Order, int, error) {
	return nil, 0, nil
}

type uosMockOrderItemRepo struct {
	items []orderDomain.OrderItem
	err   error
}

func (m *uosMockOrderItemRepo) ListByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]orderDomain.OrderItem, error) {
	return m.items, m.err
}

func (m *uosMockOrderItemRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]orderDomain.OrderItem, error) {
	return nil, nil
}

func (m *uosMockOrderItemRepo) SaveBulk(_ context.Context, _ transaction.Executor, _ []orderDomain.OrderItem) error {
	return nil
}

type uosMockShipmentRepo struct {
	createdShipment *shipmentDomain.Shipment
	err             error
}

func (m *uosMockShipmentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}

func (m *uosMockShipmentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}

func (m *uosMockShipmentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]shipmentDomain.Shipment, error) {
	return nil, nil
}

func (m *uosMockShipmentRepo) Create(_ context.Context, _ transaction.Executor, shipment shipmentDomain.Shipment) error {
	m.createdShipment = &shipment
	return m.err
}

func (m *uosMockShipmentRepo) Update(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}

type uosMockAddressRepo struct {
	addr *addressDomain.CustomerAddress
	err  error
}

func (m *uosMockAddressRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*addressDomain.CustomerAddress, error) {
	if m.addr != nil && m.addr.ID == id {
		return m.addr, nil
	}
	return nil, m.err
}

func (m *uosMockAddressRepo) GetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.CustomerAddress, error) {
	return nil, nil
}

func (m *uosMockAddressRepo) ListByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]addressDomain.CustomerAddress, error) {
	return nil, nil
}

func (m *uosMockAddressRepo) CountByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*int, error) {
	return nil, nil
}

func (m *uosMockAddressRepo) UnsetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}

func (m *uosMockAddressRepo) Save(_ context.Context, _ transaction.Executor, _ addressDomain.CustomerAddress) error {
	return nil
}

func (m *uosMockAddressRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}

func (m *uosMockAddressRepo) DeleteByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}

func (m *uosMockAddressRepo) ListByIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]addressDomain.CustomerAddress, error) {
	if m.addr != nil {
		return []addressDomain.CustomerAddress{*m.addr}, m.err
	}
	return nil, m.err
}

type uosMockShopAddressRepo struct {
	addr *addressDomain.ShopAddress
	err  error
}

func (m *uosMockShopAddressRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.ShopAddress, error) {
	return nil, nil
}

func (m *uosMockShopAddressRepo) GetDefaultByShopID(_ context.Context, _ transaction.Executor, shopID uuid.UUID) (*addressDomain.ShopAddress, error) {
	if m.addr != nil && m.addr.ShopID == shopID {
		return m.addr, nil
	}
	return nil, m.err
}

func (m *uosMockShopAddressRepo) GetDefaultsByShopIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) (map[uuid.UUID]addressDomain.ShopAddress, error) {
	return nil, nil
}

func (m *uosMockShopAddressRepo) FindByShopID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]addressDomain.ShopAddress, error) {
	return nil, nil
}

func (m *uosMockShopAddressRepo) Create(_ context.Context, _ transaction.Executor, _ addressDomain.ShopAddress) error {
	return nil
}

type uosMockLogisticsProvider struct {
	calledInput *shipping.CreateOrderInput
	result      *shipping.CreateOrderResult
	err         error
}

func (m *uosMockLogisticsProvider) CreateOrder(_ context.Context, input shipping.CreateOrderInput) (*shipping.CreateOrderResult, error) {
	m.calledInput = &input
	return m.result, m.err
}

func (m *uosMockLogisticsProvider) TrackShipment(_ context.Context, _ shipping.TrackShipmentInput) ([]shipping.TrackingEvent, error) {
	return nil, nil
}

type uosMockAuditLogger struct {
	events []applogger.AuditEvent
}

func (m *uosMockAuditLogger) Log(_ context.Context, event applogger.AuditEvent) {
	m.events = append(m.events, event)
}

// ===========================================================================
// Tests
// ===========================================================================

func TestUpdateOrderStatus_SimpleTransition(t *testing.T) {
	orderID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusConfirmed,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Order.Status != orderDomain.OrderStatusConfirmed {
		t.Errorf("expected status Confirmed, got %s", res.Order.Status)
	}

	if res.Shipment != nil {
		t.Errorf("expected no shipment for simple transition, got %+v", res.Shipment)
	}

	if orderRepo.updatedState.status != orderDomain.OrderStatusConfirmed {
		t.Errorf("expected status Confirmed in DB, got %s", orderRepo.updatedState.status)
	}
}

func TestUpdateOrderStatus_ShippedCourierManual(t *testing.T) {
	orderID := uuid.New()
	addressID := uuid.New()
	shopID := uuid.New()
	order := &orderDomain.Order{
		ID:          orderID,
		Number:      "ORD-123",
		Status:      orderDomain.OrderStatusProcessing,
		AddressID:   addressID,
		ShippingFee: 15000,
	}

	courierCode := "jne"
	courierService := "reg"
	items := []orderDomain.OrderItem{
		{
			ID:             uuid.New(),
			OrderID:        orderID,
			ShopID:         shopID,
			ProductName:    "Flower A",
			CourierCode:    &courierCode,
			CourierService: &courierService,
		},
	}

	custAddr := &addressDomain.CustomerAddress{
		ID:           addressID,
		ReceiverName: "Bob",
		Detail: addressDomain.AddressDetail{
			DistrictID:  "456",
			FullAddress: "Receiver Rd 12",
		},
	}

	shopAddr := &addressDomain.ShopAddress{
		ShopID: shopID,
		Label:  "Main Shop",
		Detail: addressDomain.AddressDetail{
			DistrictID:  "123",
			FullAddress: "Shipper Rd 99",
		},
	}

	orderRepo := &uosMockOrderRepo{order: order}
	orderItemRepo := &uosMockOrderItemRepo{items: items}
	shipmentRepo := &uosMockShipmentRepo{}
	addressRepo := &uosMockAddressRepo{addr: custAddr}
	shopAddressRepo := &uosMockShopAddressRepo{addr: shopAddr}

	manualTracking := "TRACK12345"
	logisticsMock := &uosMockLogisticsProvider{
		result: &shipping.CreateOrderResult{
			TrackingNumber: manualTracking,
		},
	}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		orderItemRepo,
		shipmentRepo,
		addressRepo,
		shopAddressRepo,
		logisticsMock,
		&uosMockAuditLogger{},
	)

	// Verify courier mode with manual tracking input (Milestone 1 check)
	fulfillmentMethod := "courier"
	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID:           orderID,
		Status:            orderDomain.OrderStatusShipped,
		TrackingNumber:    &manualTracking,
		FulfillmentMethod: &fulfillmentMethod,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Order.Status != orderDomain.OrderStatusShipped {
		t.Errorf("expected status Shipped, got %s", res.Order.Status)
	}

	if res.Shipment == nil {
		t.Fatal("expected shipment to be created")
	}

	if res.Shipment.FulfillmentMethod != shipmentDomain.FulfillmentMethodCourier {
		t.Errorf("expected courier fulfillment, got %s", res.Shipment.FulfillmentMethod)
	}

	if *res.Shipment.TrackingNumber != manualTracking {
		t.Errorf("expected tracking number %s, got %v", manualTracking, res.Shipment.TrackingNumber)
	}

	if logisticsMock.calledInput == nil {
		t.Fatal("expected logistics provider CreateOrder to be called")
	}

	if *logisticsMock.calledInput.ManualTrackingNumber != manualTracking {
		t.Errorf("expected logistics input to have manual tracking %s", manualTracking)
	}

	if res.Shipment.Courier != courierCode || res.Shipment.Service != courierService {
		t.Errorf("expected courier %s and service %s, got %s and %s", courierCode, courierService, res.Shipment.Courier, res.Shipment.Service)
	}
}

func TestUpdateOrderStatus_ShippedSelfDelivery(t *testing.T) {
	orderID := uuid.New()
	addressID := uuid.New()
	shopID := uuid.New()
	order := &orderDomain.Order{
		ID:          orderID,
		Number:      "ORD-456",
		Status:      orderDomain.OrderStatusProcessing,
		AddressID:   addressID,
		ShippingFee: 0,
	}

	items := []orderDomain.OrderItem{
		{
			ID:          uuid.New(),
			OrderID:     orderID,
			ShopID:      shopID,
			ProductName: "Flower B",
		},
	}

	custAddr := &addressDomain.CustomerAddress{
		ID:           addressID,
		ReceiverName: "Alice",
		Detail: addressDomain.AddressDetail{
			DistrictID:  "456",
			FullAddress: "Receiver Rd 12",
		},
	}

	shopAddr := &addressDomain.ShopAddress{
		ShopID: shopID,
		Label:  "Main Shop",
		Detail: addressDomain.AddressDetail{
			DistrictID:  "123",
			FullAddress: "Shipper Rd 99",
		},
	}

	orderRepo := &uosMockOrderRepo{order: order}
	orderItemRepo := &uosMockOrderItemRepo{items: items}
	shipmentRepo := &uosMockShipmentRepo{}
	addressRepo := &uosMockAddressRepo{addr: custAddr}
	shopAddressRepo := &uosMockShopAddressRepo{addr: shopAddr}

	logisticsMock := &uosMockLogisticsProvider{}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		orderItemRepo,
		shipmentRepo,
		addressRepo,
		shopAddressRepo,
		logisticsMock,
		&uosMockAuditLogger{},
	)

	// Verify Milestone 2: Self delivery
	fulfillmentMethod := "self_delivery"
	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID:           orderID,
		Status:            orderDomain.OrderStatusShipped,
		FulfillmentMethod: &fulfillmentMethod,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Order.Status != orderDomain.OrderStatusShipped {
		t.Errorf("expected status Shipped, got %s", res.Order.Status)
	}

	if res.Shipment == nil {
		t.Fatal("expected shipment to be created")
	}

	if res.Shipment.FulfillmentMethod != shipmentDomain.FulfillmentMethodSelfDelivery {
		t.Errorf("expected self_delivery fulfillment, got %s", res.Shipment.FulfillmentMethod)
	}

	if res.Shipment.TrackingNumber != nil {
		t.Errorf("expected nil tracking number for self_delivery, got %v", *res.Shipment.TrackingNumber)
	}

	if logisticsMock.calledInput != nil {
		t.Error("expected logistics provider NOT to be called for self_delivery")
	}

	if res.Shipment.Courier != "self_delivery" || res.Shipment.Service != "self_delivery" {
		t.Errorf("expected courier and service to be 'self_delivery', got %s and %s", res.Shipment.Courier, res.Shipment.Service)
	}
}

func TestUpdateOrderStatus_InvalidTransition(t *testing.T) {
	orderID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

	// Can't transition directly from Pending to Shipped
	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusShipped,
	})

	if err == nil {
		t.Fatal("expected transition error, got nil")
	}

	if res != nil {
		t.Errorf("expected nil result, got %+v", res)
	}
}

func TestUpdateOrderStatus_DefaultFulfillment(t *testing.T) {
	orderID := uuid.New()
	addressID := uuid.New()
	shopID := uuid.New()
	order := &orderDomain.Order{
		ID:          orderID,
		Number:      "ORD-789",
		Status:      orderDomain.OrderStatusProcessing,
		AddressID:   addressID,
		ShippingFee: 12000,
	}

	courierCode := "pos"
	courierService := "khusus"
	items := []orderDomain.OrderItem{
		{
			ID:             uuid.New(),
			OrderID:        orderID,
			ShopID:         shopID,
			ProductName:    "Flower C",
			CourierCode:    &courierCode,
			CourierService: &courierService,
		},
	}

	custAddr := &addressDomain.CustomerAddress{
		ID:           addressID,
		ReceiverName: "Charlie",
		Detail: addressDomain.AddressDetail{
			DistrictID:  "456",
			FullAddress: "Receiver Rd 12",
		},
	}

	shopAddr := &addressDomain.ShopAddress{
		ShopID: shopID,
		Label:  "Main Shop",
		Detail: addressDomain.AddressDetail{
			DistrictID:  "123",
			FullAddress: "Shipper Rd 99",
		},
	}

	orderRepo := &uosMockOrderRepo{order: order}
	orderItemRepo := &uosMockOrderItemRepo{items: items}
	shipmentRepo := &uosMockShipmentRepo{}
	addressRepo := &uosMockAddressRepo{addr: custAddr}
	shopAddressRepo := &uosMockShopAddressRepo{addr: shopAddr}

	logisticsMock := &uosMockLogisticsProvider{
		result: &shipping.CreateOrderResult{
			TrackingNumber: "POS9999",
		},
	}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		orderItemRepo,
		shipmentRepo,
		addressRepo,
		shopAddressRepo,
		logisticsMock,
		&uosMockAuditLogger{},
	)

	// If FulfillmentMethod is not provided, it should default to courier
	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusShipped,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Shipment.FulfillmentMethod != shipmentDomain.FulfillmentMethodCourier {
		t.Errorf("expected default fulfillment 'courier', got %s", res.Shipment.FulfillmentMethod)
	}

	if logisticsMock.calledInput == nil {
		t.Fatal("expected logistics provider to be called by default")
	}
}

func TestUpdateOrderStatus_AuditLogging_Success(t *testing.T) {
	orderID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusConfirmed,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	auditLogger := &uosMockAuditLogger{}
	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		auditLogger,
	)

	_, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusProcessing,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(auditLogger.events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(auditLogger.events))
	}

	event := auditLogger.events[0]
	if event.Category != "user_action" {
		t.Errorf("expected Category 'user_action', got %q", event.Category)
	}
	if event.Action != "update_order_status" {
		t.Errorf("expected Action 'update_order_status', got %q", event.Action)
	}
	if event.Resource != "order" {
		t.Errorf("expected Resource 'order', got %q", event.Resource)
	}
	if event.ResourceID != orderID.String() {
		t.Errorf("expected ResourceID %q, got %q", orderID.String(), event.ResourceID)
	}
	if event.Outcome != applogger.OutcomeSuccess {
		t.Errorf("expected Outcome 'success', got %q", event.Outcome)
	}
	oldStatusVal, ok := event.Metadata["old_status"].(string)
	if !ok || oldStatusVal != "confirmed" {
		t.Errorf("expected old_status 'confirmed', got %v", event.Metadata["old_status"])
	}
	newStatusVal, ok := event.Metadata["new_status"].(string)
	if !ok || newStatusVal != "processing" {
		t.Errorf("expected new_status 'processing', got %v", event.Metadata["new_status"])
	}
}

func TestUpdateOrderStatus_AuditLogging_Failure(t *testing.T) {
	orderID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	auditLogger := &uosMockAuditLogger{}
	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		auditLogger,
	)

	// Invalid status transition Pending -> Delivered
	_, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusDelivered,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if len(auditLogger.events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(auditLogger.events))
	}

	event := auditLogger.events[0]
	if event.Outcome != applogger.OutcomeFailure {
		t.Errorf("expected Outcome 'failure', got %q", event.Outcome)
	}
	if event.Metadata["error"] == nil {
		t.Errorf("expected error in metadata, got nil")
	}
}
