package usecase

import (
	"context"
	"testing"
	"time"

	applogger "service-core/internal/common/logger"
	shipping "service-core/internal/infra/shipping"
	addressDomain "service-core/internal/modules/address/domain"
	inventoryDomain "service-core/internal/modules/inventory/domain"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	paymentDomain "service-core/internal/modules/payment/domain"
	productDomain "service-core/internal/modules/product/domain"
	productRepo "service-core/internal/modules/product/repository"
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

func (m *uosMockOrderRepo) UpdateStatusWithSLA(_ context.Context, _ transaction.Executor, id uuid.UUID, status orderDomain.OrderStatus, _ *time.Time, _ *time.Time) error {
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
func (m *uosMockOrderRepo) SetConfirmedAndExpiry(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ time.Time, _ time.Time) error {
	return nil
}
func (m *uosMockOrderRepo) FindExpiredUnfulfilledOrders(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]orderDomain.Order, error) {
	return nil, nil
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

type uosMockInventoryRepo struct {
	committed map[uuid.UUID]int
	released  map[uuid.UUID]int
	err       error
}

func (m *uosMockInventoryRepo) GetByProductIDAndShopID(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ uuid.UUID) (*inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *uosMockInventoryRepo) ListByProductID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *uosMockInventoryRepo) ListByProductIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) (map[uuid.UUID][]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *uosMockInventoryRepo) ListByShopID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]inventoryDomain.Inventory, error) {
	return nil, nil
}
func (m *uosMockInventoryRepo) Create(_ context.Context, _ transaction.Executor, _ *inventoryDomain.Inventory) error {
	return nil
}
func (m *uosMockInventoryRepo) Reserve(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ uuid.UUID, _ int) error {
	return nil
}
func (m *uosMockInventoryRepo) Release(_ context.Context, _ transaction.Executor, productID uuid.UUID, _ uuid.UUID, qty int) error {
	if m.released == nil {
		m.released = make(map[uuid.UUID]int)
	}
	m.released[productID] += qty
	return m.err
}
func (m *uosMockInventoryRepo) Commit(_ context.Context, _ transaction.Executor, productID uuid.UUID, _ uuid.UUID, qty int) error {
	if m.committed == nil {
		m.committed = make(map[uuid.UUID]int)
	}
	m.committed[productID] += qty
	return m.err
}
func (m *uosMockInventoryRepo) Restock(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ uuid.UUID, _ int) error {
	return nil
}
func (m *uosMockInventoryRepo) Update(_ context.Context, _ transaction.Executor, _ *inventoryDomain.Inventory) error {
	return nil
}
func (m *uosMockInventoryRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ uuid.UUID) error {
	return nil
}

type uosMockPaymentRepo struct {
	payment *paymentDomain.Payment
	err     error
}

func (m *uosMockPaymentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.Payment, error) {
	return nil, nil
}
func (m *uosMockPaymentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) (*paymentDomain.Payment, error) {
	if m.payment != nil && m.payment.OrderID == orderID {
		return m.payment, nil
	}
	return nil, m.err
}
func (m *uosMockPaymentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *uosMockPaymentRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ paymentDomain.PaymentStatus) error {
	return nil
}
func (m *uosMockPaymentRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.Payment) error {
	return nil
}
func (m *uosMockPaymentRepo) ListPendingGateway(_ context.Context, _ transaction.Executor, _ time.Time) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *uosMockPaymentRepo) ListPastDuePending(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]paymentDomain.Payment, error) {
	return nil, nil
}

type uosMockShipmentRepo struct {
	shipment         *shipmentDomain.Shipment
	createdShipment  *shipmentDomain.Shipment
	createdShipments []shipmentDomain.Shipment
	updatedShipment  *shipmentDomain.Shipment
	err              error
}

func (m *uosMockShipmentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}

func (m *uosMockShipmentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) (*shipmentDomain.Shipment, error) {
	if m.shipment != nil && m.shipment.OrderID == orderID {
		return m.shipment, nil
	}
	return nil, m.err
}

func (m *uosMockShipmentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]shipmentDomain.Shipment, error) {
	return nil, nil
}

func (m *uosMockShipmentRepo) Create(_ context.Context, _ transaction.Executor, shipment shipmentDomain.Shipment) error {
	m.createdShipment = &shipment
	m.createdShipments = append(m.createdShipments, shipment)
	return m.err
}

func (m *uosMockShipmentRepo) Update(_ context.Context, _ transaction.Executor, shipment shipmentDomain.Shipment) error {
	m.updatedShipment = &shipment
	return m.err
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
	addr  *addressDomain.ShopAddress
	addrs map[uuid.UUID]*addressDomain.ShopAddress
	err   error
}

func (m *uosMockShopAddressRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.ShopAddress, error) {
	return nil, nil
}

func (m *uosMockShopAddressRepo) GetDefaultByShopID(_ context.Context, _ transaction.Executor, shopID uuid.UUID) (*addressDomain.ShopAddress, error) {
	if m.addrs != nil {
		if a, ok := m.addrs[shopID]; ok {
			return a, nil
		}
	}
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

type uosMockProductRepo struct {
	products map[uuid.UUID]productDomain.Product
	err      error
}

func (m *uosMockProductRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*productDomain.Product, error) {
	if p, ok := m.products[id]; ok {
		return &p, nil
	}
	return nil, m.err
}

func (m *uosMockProductRepo) GetBySlug(_ context.Context, _ transaction.Executor, _ string) (*productDomain.Product, error) {
	return nil, nil
}

func (m *uosMockProductRepo) FindProducts(_ context.Context, _ transaction.Executor, _ productRepo.FindProductParams) ([]productDomain.Product, int, error) {
	return nil, 0, nil
}

func (m *uosMockProductRepo) FindProductsWithInventory(_ context.Context, _ transaction.Executor, _ productRepo.FindProductParams) ([]productDomain.ProductWithInventory, int, error) {
	return nil, 0, nil
}

func (m *uosMockProductRepo) FindByIDs(_ context.Context, _ transaction.Executor, ids []uuid.UUID) ([]productDomain.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []productDomain.Product
	for _, id := range ids {
		if p, ok := m.products[id]; ok {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *uosMockProductRepo) Save(_ context.Context, _ transaction.Executor, _ *productDomain.Product) error {
	return nil
}

func (m *uosMockProductRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
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

func TestUpdateOrderStatus_Confirmed_PaymentNotPaid(t *testing.T) {
	orderID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	paymentRepo := &uosMockPaymentRepo{
		payment: &paymentDomain.Payment{
			OrderID: orderID,
			Status:  paymentDomain.PaymentStatusPending,
		},
	}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockInventoryRepo{},
		paymentRepo,
		&uosMockProductRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

	_, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusConfirmed,
	})

	if err == nil {
		t.Fatal("expected error when confirming order without paid payment, got nil")
	}
}

func TestUpdateOrderStatus_Confirmed_Success(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	items := []orderDomain.OrderItem{
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			ShopID:    shopID,
			ProductID: productID,
			Quantity:  2,
		},
	}

	orderRepo := &uosMockOrderRepo{order: order}
	orderItemRepo := &uosMockOrderItemRepo{items: items}
	inventoryRepo := &uosMockInventoryRepo{}
	paymentRepo := &uosMockPaymentRepo{
		payment: &paymentDomain.Payment{
			OrderID: orderID,
			Status:  paymentDomain.PaymentStatusPaid,
		},
	}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		orderItemRepo,
		inventoryRepo,
		paymentRepo,
		&uosMockProductRepo{},
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

	if inventoryRepo.committed[productID] != 2 {
		t.Errorf("expected 2 committed for product %s, got %d", productID, inventoryRepo.committed[productID])
	}
}

func TestUpdateOrderStatus_Processing_PaymentNotPaid(t *testing.T) {
	orderID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusConfirmed,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	paymentRepo := &uosMockPaymentRepo{
		payment: &paymentDomain.Payment{
			OrderID: orderID,
			Status:  paymentDomain.PaymentStatusPending,
		},
	}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockInventoryRepo{},
		paymentRepo,
		&uosMockProductRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

	_, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusProcessing,
	})

	if err == nil {
		t.Fatal("expected error when moving to processing without paid payment, got nil")
	}
}

func TestUpdateOrderStatus_Processing_Success(t *testing.T) {
	orderID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusConfirmed,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	paymentRepo := &uosMockPaymentRepo{
		payment: &paymentDomain.Payment{
			OrderID: orderID,
			Status:  paymentDomain.PaymentStatusPaid,
		},
	}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockInventoryRepo{},
		paymentRepo,
		&uosMockProductRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusProcessing,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Order.Status != orderDomain.OrderStatusProcessing {
		t.Errorf("expected status Processing, got %s", res.Order.Status)
	}
}

func TestUpdateOrderStatus_Cancelled_ReleasesInventory(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusPending,
	}

	items := []orderDomain.OrderItem{
		{
			ID:        uuid.New(),
			OrderID:   orderID,
			ShopID:    shopID,
			ProductID: productID,
			Quantity:  3,
		},
	}

	orderRepo := &uosMockOrderRepo{order: order}
	orderItemRepo := &uosMockOrderItemRepo{items: items}
	inventoryRepo := &uosMockInventoryRepo{}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		orderItemRepo,
		inventoryRepo,
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusCancelled,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Order.Status != orderDomain.OrderStatusCancelled {
		t.Errorf("expected status Cancelled, got %s", res.Order.Status)
	}

	if inventoryRepo.released[productID] != 3 {
		t.Errorf("expected 3 released for product %s, got %d", productID, inventoryRepo.released[productID])
	}
}

func TestUpdateOrderStatus_Delivered_ClosesShipment(t *testing.T) {
	orderID := uuid.New()
	shipmentID := uuid.New()
	order := &orderDomain.Order{
		ID:     orderID,
		Status: orderDomain.OrderStatusShipped,
	}

	shipment := &shipmentDomain.Shipment{
		ID:      shipmentID,
		OrderID: orderID,
		Status:  shipmentDomain.ShipmentStatusInTransit,
	}

	orderRepo := &uosMockOrderRepo{order: order}
	shipmentRepo := &uosMockShipmentRepo{shipment: shipment}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		shipmentRepo,
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusDelivered,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Order.Status != orderDomain.OrderStatusDelivered {
		t.Errorf("expected status Delivered, got %s", res.Order.Status)
	}

	if shipmentRepo.updatedShipment == nil || shipmentRepo.updatedShipment.Status != shipmentDomain.ShipmentStatusDelivered {
		t.Errorf("expected shipment status Delivered, got %v", shipmentRepo.updatedShipment)
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
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		shipmentRepo,
		addressRepo,
		shopAddressRepo,
		logisticsMock,
		&uosMockAuditLogger{},
	)

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
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		shipmentRepo,
		addressRepo,
		shopAddressRepo,
		logisticsMock,
		&uosMockAuditLogger{},
	)

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
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		&uosMockAuditLogger{},
	)

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
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		shipmentRepo,
		addressRepo,
		shopAddressRepo,
		logisticsMock,
		&uosMockAuditLogger{},
	)

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
	paymentRepo := &uosMockPaymentRepo{
		payment: &paymentDomain.Payment{
			OrderID: orderID,
			Status:  paymentDomain.PaymentStatusPaid,
		},
	}
	auditLogger := &uosMockAuditLogger{}
	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		orderRepo,
		&uosMockOrderItemRepo{},
		&uosMockInventoryRepo{},
		paymentRepo,
		&uosMockProductRepo{},
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
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{},
		&uosMockShopAddressRepo{},
		&uosMockLogisticsProvider{},
		auditLogger,
	)

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

func TestUpdateOrderStatus_Shipped_ProductWeightCalculation(t *testing.T) {
	orderID := uuid.New()
	addressID := uuid.New()
	shopID := uuid.New()
	productID := uuid.New()

	order := &orderDomain.Order{
		ID:          orderID,
		Number:      "ORD-WEIGHT-1",
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
			ProductID:      productID,
			ProductName:    "Heavy Flower Box",
			Quantity:       2,
			CourierCode:    &courierCode,
			CourierService: &courierService,
		},
	}

	weightVal := 2500.0 // 2500g per unit, 2 units = 5000g
	products := map[uuid.UUID]productDomain.Product{
		productID: {
			ID:     productID,
			Name:   "Heavy Flower Box",
			Weight: &weightVal,
		},
	}

	custAddr := &addressDomain.CustomerAddress{
		ID:           addressID,
		ReceiverName: "David",
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

	logisticsMock := &uosMockLogisticsProvider{
		result: &shipping.CreateOrderResult{
			TrackingNumber: "WEIGHT123",
		},
	}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		&uosMockOrderRepo{order: order},
		&uosMockOrderItemRepo{items: items},
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{products: products},
		&uosMockShipmentRepo{},
		&uosMockAddressRepo{addr: custAddr},
		&uosMockShopAddressRepo{addr: shopAddr},
		logisticsMock,
		&uosMockAuditLogger{},
	)

	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusShipped,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Shipment.Weight != 5000 {
		t.Errorf("expected shipment weight 5000g, got %d", res.Shipment.Weight)
	}

	if logisticsMock.calledInput == nil || logisticsMock.calledInput.Weight != 5000 {
		t.Errorf("expected logistics input weight 5000g, got %v", logisticsMock.calledInput)
	}
}

func TestUpdateOrderStatus_Shipped_MultiShopShipment(t *testing.T) {
	orderID := uuid.New()
	addressID := uuid.New()
	shopID1 := uuid.New()
	shopID2 := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()

	order := &orderDomain.Order{
		ID:          orderID,
		Number:      "ORD-MULTI-1",
		Status:      orderDomain.OrderStatusProcessing,
		AddressID:   addressID,
		ShippingFee: 30000,
	}

	courier1 := "jne"
	service1 := "reg"
	courier2 := "sicepat"
	service2 := "gokil"

	items := []orderDomain.OrderItem{
		{
			ID:             uuid.New(),
			OrderID:        orderID,
			ShopID:         shopID1,
			ProductID:      productID1,
			ProductName:    "Flower Shop 1",
			Quantity:       1,
			ShippingFee:    15000,
			CourierCode:    &courier1,
			CourierService: &service1,
		},
		{
			ID:             uuid.New(),
			OrderID:        orderID,
			ShopID:         shopID2,
			ProductID:      productID2,
			ProductName:    "Flower Shop 2",
			Quantity:       1,
			ShippingFee:    15000,
			CourierCode:    &courier2,
			CourierService: &service2,
		},
	}

	custAddr := &addressDomain.CustomerAddress{
		ID:           addressID,
		ReceiverName: "Eva",
		Detail: addressDomain.AddressDetail{
			DistrictID:  "456",
			FullAddress: "Receiver Rd 12",
		},
	}

	shopAddr1 := &addressDomain.ShopAddress{
		ShopID: shopID1,
		Label:  "Shop One",
		Detail: addressDomain.AddressDetail{DistrictID: "101", FullAddress: "Shop 1 Rd"},
	}
	shopAddr2 := &addressDomain.ShopAddress{
		ShopID: shopID2,
		Label:  "Shop Two",
		Detail: addressDomain.AddressDetail{DistrictID: "102", FullAddress: "Shop 2 Rd"},
	}

	shopAddressRepo := &uosMockShopAddressRepo{
		addrs: map[uuid.UUID]*addressDomain.ShopAddress{
			shopID1: shopAddr1,
			shopID2: shopAddr2,
		},
	}

	logisticsMock := &uosMockLogisticsProvider{
		result: &shipping.CreateOrderResult{
			TrackingNumber: "MULTI123",
		},
	}

	shipmentRepo := &uosMockShipmentRepo{}

	usecase := NewUpdateOrderStatusUsecase(
		&uosMockExecutor{},
		&uosMockTransactor{},
		&uosMockOrderRepo{order: order},
		&uosMockOrderItemRepo{items: items},
		&uosMockInventoryRepo{},
		&uosMockPaymentRepo{},
		&uosMockProductRepo{},
		shipmentRepo,
		&uosMockAddressRepo{addr: custAddr},
		shopAddressRepo,
		logisticsMock,
		&uosMockAuditLogger{},
	)

	res, err := usecase.Execute(context.Background(), UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  orderDomain.OrderStatusShipped,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Order.Status != orderDomain.OrderStatusShipped {
		t.Errorf("expected status Shipped, got %s", res.Order.Status)
	}

	if len(shipmentRepo.createdShipments) != 2 {
		t.Fatalf("expected 2 shipments created for multi-shop order, got %d", len(shipmentRepo.createdShipments))
	}

	if shipmentRepo.createdShipments[0].Courier != courier1 || shipmentRepo.createdShipments[1].Courier != courier2 {
		t.Errorf("expected courier %s and %s, got %s and %s", courier1, courier2, shipmentRepo.createdShipments[0].Courier, shipmentRepo.createdShipments[1].Courier)
	}
}
