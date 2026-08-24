package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

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
// Mocks for GetOrderTrackingUsecase tests
// ===========================================================================

type gotMockExecutor struct{}

func (m *gotMockExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *gotMockExecutor) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *gotMockExecutor) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return nil
}


type gotMockOrderRepo struct {
	order *orderDomain.Order
	err   error
}

func (m *gotMockOrderRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*orderDomain.Order, error) {
	if m.order != nil && m.order.ID == id {
		return m.order, nil
	}
	return nil, m.err
}
func (m *gotMockOrderRepo) GetByNumber(_ context.Context, _ transaction.Executor, _ string) (*orderDomain.Order, error) {
	return nil, nil
}
func (m *gotMockOrderRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ orderDomain.OrderStatus) error {
	return nil
}
func (m *gotMockOrderRepo) UpdateStatusWithSLA(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ orderDomain.OrderStatus, _ *time.Time, _ *time.Time) error {
	return nil
}
func (m *gotMockOrderRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.Order) error {
	return nil
}
func (m *gotMockOrderRepo) FindOrders(_ context.Context, _ transaction.Executor, _ orderRepo.FindOrderParams) ([]orderDomain.Order, int, error) {
	return nil, 0, nil
}
func (m *gotMockOrderRepo) SetConfirmedAndExpiry(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ time.Time, _ time.Time) error {
	return nil
}
func (m *gotMockOrderRepo) FindExpiredUnfulfilledOrders(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]orderDomain.Order, error) {
	return nil, nil
}

type gotMockShipmentRepo struct {
	shipment *shipmentDomain.Shipment
	err      error
}

func (m *gotMockShipmentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *gotMockShipmentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) (*shipmentDomain.Shipment, error) {
	if m.shipment != nil && m.shipment.OrderID == orderID {
		return m.shipment, nil
	}
	return nil, m.err
}

func (m *gotMockShipmentRepo) ListByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) ([]shipmentDomain.Shipment, error) {
	if m.shipment != nil && m.shipment.OrderID == orderID {
		return []shipmentDomain.Shipment{*m.shipment}, nil
	}
	return nil, m.err
}

func (m *gotMockShipmentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *gotMockShipmentRepo) Create(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}
func (m *gotMockShipmentRepo) Update(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}

type gotMockShipmentEventRepo struct {
	events []shipmentDomain.ShipmentEvent
	err    error
}

func (m *gotMockShipmentEventRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.ShipmentEvent, error) {
	return nil, nil
}
func (m *gotMockShipmentEventRepo) ListByShipmentID(_ context.Context, _ transaction.Executor, shipmentID uuid.UUID) ([]shipmentDomain.ShipmentEvent, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}
func (m *gotMockShipmentEventRepo) Create(_ context.Context, _ transaction.Executor, _ shipmentDomain.ShipmentEvent) error {
	return nil
}

type gotMockAddressRepo struct {
	addr *addressDomain.CustomerAddress
}

func (m *gotMockAddressRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*addressDomain.CustomerAddress, error) {
	if m.addr != nil && m.addr.ID == id {
		return m.addr, nil
	}
	return nil, nil
}
func (m *gotMockAddressRepo) GetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *gotMockAddressRepo) ListByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]addressDomain.CustomerAddress, error) {
	return nil, nil
}
func (m *gotMockAddressRepo) CountByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*int, error) {
	return nil, nil
}
func (m *gotMockAddressRepo) UnsetDefaultByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *gotMockAddressRepo) Save(_ context.Context, _ transaction.Executor, _ addressDomain.CustomerAddress) error {
	return nil
}
func (m *gotMockAddressRepo) Delete(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *gotMockAddressRepo) DeleteByCustomerID(_ context.Context, _ transaction.Executor, _ uuid.UUID) error {
	return nil
}
func (m *gotMockAddressRepo) ListByIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]addressDomain.CustomerAddress, error) {
	if m.addr != nil {
		return []addressDomain.CustomerAddress{*m.addr}, nil
	}
	return nil, nil
}

type gotMockLogisticsProvider struct {
	events []shipping.TrackingEvent
	err    error
}

func (m *gotMockLogisticsProvider) CreateOrder(_ context.Context, _ shipping.CreateOrderInput) (*shipping.CreateOrderResult, error) {
	return nil, nil
}
func (m *gotMockLogisticsProvider) CancelOrder(_ context.Context, _ string) error {
	return nil
}
func (m *gotMockLogisticsProvider) TrackShipment(_ context.Context, _ shipping.TrackShipmentInput) ([]shipping.TrackingEvent, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

// ===========================================================================
// Test Cases
// ===========================================================================

func TestGetOrderTracking_OrderNotFound(t *testing.T) {
	usecase := NewGetOrderTrackingUsecase(
		&gotMockExecutor{},
		&gotMockOrderRepo{order: nil},
		&gotMockShipmentRepo{},
		&gotMockShipmentEventRepo{},
		&gotMockLogisticsProvider{},
		&gotMockAddressRepo{},
		nil,
	)

	result, err := usecase.Execute(context.Background(), GetOrderTrackingInput{
		OrderID:    uuid.New(),
		CustomerID: uuid.New(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "order not found" {
		t.Fatalf("expected 'order not found' error, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestGetOrderTracking_CustomerIDMismatch(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()
	mismatchCustomerID := uuid.New()

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
	}

	usecase := NewGetOrderTrackingUsecase(
		&gotMockExecutor{},
		&gotMockOrderRepo{order: order},
		&gotMockShipmentRepo{},
		&gotMockShipmentEventRepo{},
		&gotMockLogisticsProvider{},
		&gotMockAddressRepo{},
		nil,
	)

	result, err := usecase.Execute(context.Background(), GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: mismatchCustomerID,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "not authorized" {
		t.Fatalf("expected 'not authorized' error, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestGetOrderTracking_ShipmentNotFound(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
	}

	usecase := NewGetOrderTrackingUsecase(
		&gotMockExecutor{},
		&gotMockOrderRepo{order: order},
		&gotMockShipmentRepo{shipment: nil},
		&gotMockShipmentEventRepo{},
		&gotMockLogisticsProvider{},
		&gotMockAddressRepo{},
		nil,
	)

	result, err := usecase.Execute(context.Background(), GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: customerID,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "shipment not found" {
		t.Fatalf("expected 'shipment not found' error, got: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
}

func TestGetOrderTracking_InternalEventsOnly(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()
	shipmentID := uuid.New()

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
	}

	trackingNo := "TRACK123"
	shipment := &shipmentDomain.Shipment{
		ID:                shipmentID,
		OrderID:           orderID,
		FulfillmentMethod: shipmentDomain.FulfillmentMethodCourier,
		TrackingNumber:    &trackingNo,
		Courier:           "jne",
	}

	now := time.Now()
	internalEvents := []shipmentDomain.ShipmentEvent{
		{
			ID:          uuid.New(),
			ShipmentID:  shipmentID,
			Status:      "labelled",
			Description: "Label generated",
			Timestamp:   now.Add(1 * time.Hour),
		},
		{
			ID:          uuid.New(),
			ShipmentID:  shipmentID,
			Status:      "packed",
			Description: "Items packed",
			Timestamp:   now,
		},
	}

	usecase := NewGetOrderTrackingUsecase(
		&gotMockExecutor{},
		&gotMockOrderRepo{order: order},
		&gotMockShipmentRepo{shipment: shipment},
		&gotMockShipmentEventRepo{events: internalEvents},
		&gotMockLogisticsProvider{events: nil, err: nil}, // returns no external events
		&gotMockAddressRepo{},
		nil,
	)

	result, err := usecase.Execute(context.Background(), GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: customerID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.OrderID != orderID || result.ShipmentID != shipmentID {
		t.Fatalf("unexpected order/shipment ID")
	}

	if len(result.Timeline) != 2 {
		t.Fatalf("expected 2 timeline events, got %d", len(result.Timeline))
	}

	// Verify chronological order (Timestamp ASC)
	if result.Timeline[0].Status != "packed" || result.Timeline[1].Status != "labelled" {
		t.Fatalf("timeline events not sorted chronologically: first=%s, second=%s", result.Timeline[0].Status, result.Timeline[1].Status)
	}
}

func TestGetOrderTracking_MergedEvents(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()
	shipmentID := uuid.New()
	addressID := uuid.New()

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
		AddressID:  addressID,
	}

	trackingNo := "TRACK123"
	shipment := &shipmentDomain.Shipment{
		ID:                shipmentID,
		OrderID:           orderID,
		FulfillmentMethod: shipmentDomain.FulfillmentMethodCourier,
		TrackingNumber:    &trackingNo,
		Courier:           "jne",
	}

	now := time.Now()
	internalEvents := []shipmentDomain.ShipmentEvent{
		{
			ID:          uuid.New(),
			ShipmentID:  shipmentID,
			Status:      "packed",
			Description: "Items packed",
			Timestamp:   now,
		},
		{
			ID:          uuid.New(),
			ShipmentID:  shipmentID,
			Status:      "picked_up",
			Description: "Picked up by courier",
			Timestamp:   now.Add(2 * time.Hour),
		},
	}

	externalEvents := []shipping.TrackingEvent{
		{
			Status:      "transit",
			Description: "Departed JNE Hub",
			Location:    "Jakarta",
			Timestamp:   now.Add(3 * time.Hour),
		},
		{
			Status:      "manifested",
			Description: "Shipment booked",
			Location:    "Bandung",
			Timestamp:   now.Add(1 * time.Hour),
		},
	}

	phone := "08123456789"
	addr := &addressDomain.CustomerAddress{
		ID:    addressID,
		Phone: &phone,
	}

	usecase := NewGetOrderTrackingUsecase(
		&gotMockExecutor{},
		&gotMockOrderRepo{order: order},
		&gotMockShipmentRepo{shipment: shipment},
		&gotMockShipmentEventRepo{events: internalEvents},
		&gotMockLogisticsProvider{events: externalEvents},
		&gotMockAddressRepo{addr: addr},
		nil,
	)

	result, err := usecase.Execute(context.Background(), GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: customerID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Timeline) != 4 {
		t.Fatalf("expected 4 combined timeline events, got %d", len(result.Timeline))
	}

	// Verify sorting (Timestamp ASC)
	// 0: packed (now)
	// 1: manifested (now + 1h)
	// 2: picked_up (now + 2h)
	// 3: transit (now + 3h)
	expectedOrder := []string{"packed", "manifested", "picked_up", "transit"}
	for i, status := range expectedOrder {
		if result.Timeline[i].Status != status {
			t.Errorf("at index %d: expected status %q, got %q", i, status, result.Timeline[i].Status)
		}
	}
}

func TestGetOrderTracking_ExternalAPIErrorResilience(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()
	shipmentID := uuid.New()

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
	}

	trackingNo := "TRACK123"
	shipment := &shipmentDomain.Shipment{
		ID:                shipmentID,
		OrderID:           orderID,
		FulfillmentMethod: shipmentDomain.FulfillmentMethodCourier,
		TrackingNumber:    &trackingNo,
		Courier:           "jne",
	}

	now := time.Now()
	internalEvents := []shipmentDomain.ShipmentEvent{
		{
			ID:          uuid.New(),
			ShipmentID:  shipmentID,
			Status:      "packed",
			Description: "Items packed",
			Timestamp:   now,
		},
	}

	usecase := NewGetOrderTrackingUsecase(
		&gotMockExecutor{},
		&gotMockOrderRepo{order: order},
		&gotMockShipmentRepo{shipment: shipment},
		&gotMockShipmentEventRepo{events: internalEvents},
		&gotMockLogisticsProvider{err: errors.New("external provider timeout")},
		&gotMockAddressRepo{},
		nil,
	)

	result, err := usecase.Execute(context.Background(), GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: customerID,
	})

	// Usecase should succeed and ignore external API failures
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Timeline) != 1 {
		t.Fatalf("expected 1 timeline event, got %d", len(result.Timeline))
	}
	if result.Timeline[0].Status != "packed" {
		t.Fatalf("expected packed status, got %q", result.Timeline[0].Status)
	}
}

func TestGetOrderTracking_CacheHit(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()
	shipmentID := uuid.New()
	trackingNo := "TRACK_CACHE_123"

	order := &orderDomain.Order{
		ID:         orderID,
		CustomerID: customerID,
	}

	shipment := &shipmentDomain.Shipment{
		ID:                shipmentID,
		OrderID:           orderID,
		FulfillmentMethod: shipmentDomain.FulfillmentMethodCourier,
		TrackingNumber:    &trackingNo,
		Courier:           "jne",
	}

	cache := shipping.NewTrackingCache(5 * time.Minute)
	initialEvents := []shipping.TrackingEvent{
		{
			Status:      "in_transit",
			Description: "On the way to hub",
			Location:    "JAKARTA",
			Timestamp:   time.Now(),
		},
	}
	cache.Set("jne", trackingNo, initialEvents)

	// Logistics provider is configured to return error if called
	logisticsMock := &gotMockLogisticsProvider{err: errors.New("should not be called due to cache hit")}

	usecase := NewGetOrderTrackingUsecase(
		&gotMockExecutor{},
		&gotMockOrderRepo{order: order},
		&gotMockShipmentRepo{shipment: shipment},
		&gotMockShipmentEventRepo{},
		logisticsMock,
		&gotMockAddressRepo{},
		cache,
	)

	result, err := usecase.Execute(context.Background(), GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: customerID,
	})

	if err != nil {
		t.Fatalf("expected no error on cache hit, got %v", err)
	}
	if result == nil || len(result.Timeline) != 1 {
		t.Fatalf("expected 1 cached timeline event, got %v", result)
	}
	if result.Timeline[0].Status != "in_transit" {
		t.Errorf("expected status 'in_transit', got %q", result.Timeline[0].Status)
	}
}
