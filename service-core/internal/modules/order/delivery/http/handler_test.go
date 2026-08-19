package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apperrors "service-core/internal/common/errors"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	orderDomain "service-core/internal/modules/order/domain"
	orderRepo "service-core/internal/modules/order/repository"
	"service-core/internal/modules/order/usecase"
	paymentDomain "service-core/internal/modules/payment/domain"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	transaction "service-core/internal/shared/transaction"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Mocks for handler tests
type hMockExecutor struct{}

func (m *hMockExecutor) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}
func (m *hMockExecutor) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, nil
}
func (m *hMockExecutor) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row {
	return nil
}

type hMockOrderRepo struct {
	order *orderDomain.Order
}

func (m *hMockOrderRepo) GetByID(_ context.Context, _ transaction.Executor, id uuid.UUID) (*orderDomain.Order, error) {
	if m.order != nil && m.order.ID == id {
		return m.order, nil
	}
	return nil, nil
}
func (m *hMockOrderRepo) GetByNumber(_ context.Context, _ transaction.Executor, _ string) (*orderDomain.Order, error) {
	return nil, nil
}
func (m *hMockOrderRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ orderDomain.OrderStatus) error {
	return nil
}
func (m *hMockOrderRepo) UpdateStatusWithSLA(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ orderDomain.OrderStatus, _ *time.Time, _ *time.Time) error {
	return nil
}
func (m *hMockOrderRepo) Save(_ context.Context, _ transaction.Executor, _ orderDomain.Order) error {
	return nil
}
func (m *hMockOrderRepo) FindOrders(_ context.Context, _ transaction.Executor, _ orderRepo.FindOrderParams) ([]orderDomain.Order, int, error) {
	return nil, 0, nil
}
func (m *hMockOrderRepo) SetConfirmedAndExpiry(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ time.Time, _ time.Time) error {
	return nil
}
func (m *hMockOrderRepo) FindExpiredUnfulfilledOrders(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]orderDomain.Order, error) {
	return nil, nil
}

type hMockOrderItemRepo struct {
	items []orderDomain.OrderItem
}

func (m *hMockOrderItemRepo) ListByOrderID(_ context.Context, _ transaction.Executor, orderID uuid.UUID) ([]orderDomain.OrderItem, error) {
	return m.items, nil
}
func (m *hMockOrderItemRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]orderDomain.OrderItem, error) {
	return m.items, nil
}
func (m *hMockOrderItemRepo) ListByShipmentID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]orderDomain.OrderItem, error) {
	return nil, nil
}
func (m *hMockOrderItemRepo) SaveBulk(_ context.Context, _ transaction.Executor, _ []orderDomain.OrderItem) error {
	return nil
}
func (m *hMockOrderItemRepo) AssignShipment(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ []uuid.UUID) error {
	return nil
}

type hMockPaymentRepo struct{}

func (m *hMockPaymentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.Payment, error) {
	return nil, nil
}
func (m *hMockPaymentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.Payment, error) {
	return nil, nil
}
func (m *hMockPaymentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *hMockPaymentRepo) ListPendingGateway(_ context.Context, _ transaction.Executor, _ time.Time) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *hMockPaymentRepo) ListPastDuePending(_ context.Context, _ transaction.Executor, _ time.Time, _ int) ([]paymentDomain.Payment, error) {
	return nil, nil
}
func (m *hMockPaymentRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.Payment) error {
	return nil
}
func (m *hMockPaymentRepo) UpdateStatus(_ context.Context, _ transaction.Executor, _ uuid.UUID, _ paymentDomain.PaymentStatus) error {
	return nil
}

type hMockPaymentChannelDataRepo struct{}

func (m *hMockPaymentChannelDataRepo) GetByPaymentID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*paymentDomain.PaymentChannelData, error) {
	return nil, nil
}
func (m *hMockPaymentChannelDataRepo) ListByPaymentIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) (map[uuid.UUID]*paymentDomain.PaymentChannelData, error) {
	return nil, nil
}
func (m *hMockPaymentChannelDataRepo) Save(_ context.Context, _ transaction.Executor, _ paymentDomain.PaymentChannelData) error {
	return nil
}

type hMockShipmentRepo struct{}

func (m *hMockShipmentRepo) Create(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}
func (m *hMockShipmentRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *hMockShipmentRepo) GetByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *hMockShipmentRepo) ListByOrderID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *hMockShipmentRepo) ListByOrderIDs(_ context.Context, _ transaction.Executor, _ []uuid.UUID) ([]shipmentDomain.Shipment, error) {
	return nil, nil
}
func (m *hMockShipmentRepo) Save(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}
func (m *hMockShipmentRepo) Update(_ context.Context, _ transaction.Executor, _ shipmentDomain.Shipment) error {
	return nil
}

type hMockShipmentEventRepo struct{}

func (m *hMockShipmentEventRepo) Create(_ context.Context, _ transaction.Executor, _ shipmentDomain.ShipmentEvent) error {
	return nil
}
func (m *hMockShipmentEventRepo) GetByID(_ context.Context, _ transaction.Executor, _ uuid.UUID) (*shipmentDomain.ShipmentEvent, error) {
	return nil, nil
}
func (m *hMockShipmentEventRepo) ListByShipmentID(_ context.Context, _ transaction.Executor, _ uuid.UUID) ([]shipmentDomain.ShipmentEvent, error) {
	return nil, nil
}

// Tests

func TestResolveShopFilter_NilGetShopDependency_ReturnsError(t *testing.T) {
	handler := &orderHandler{
		getShop: nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/orders?shop_slug=test-shop", nil)
	_, specified, err := handler.resolveShopFilter(req)

	if !specified {
		t.Errorf("expected specified=true when shop_slug query param is passed")
	}
	if err == nil {
		t.Fatalf("expected internal server error when getShop dependency is nil, got nil")
	}
	appErr, ok := err.(*apperrors.AppError)
	if !ok || appErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected HTTP status 500, got %v", err)
	}
}

func TestGetOrderTrackingForStaff_MissingPermission_ReturnsForbidden(t *testing.T) {
	staffID := uuid.New()
	shopID := uuid.New()
	unassignedShopID := uuid.New()
	orderID := uuid.New()

	order := &orderDomain.Order{
		ID: orderID,
	}
	items := []orderDomain.OrderItem{
		{ID: uuid.New(), OrderID: orderID, ShopID: unassignedShopID},
	}

	orderRepoMock := &hMockOrderRepo{order: order}
	orderItemRepoMock := &hMockOrderItemRepo{items: items}

	getOrderUC := usecase.NewGetOrderUsecase(
		&hMockExecutor{},
		orderRepoMock,
		orderItemRepoMock,
		&hMockPaymentRepo{},
		&hMockPaymentChannelDataRepo{},
		&hMockShipmentRepo{},
		&hMockShipmentEventRepo{},
	)

	handler := &orderHandler{
		getOrder: getOrderUC,
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String()+"/tracking", nil)
	actor := &authzDomain.Actor{
		Type:    authenDomain.AccountTypeStaff,
		StaffID: &staffID,
		Permissions: map[uuid.UUID][]string{
			shopID: {authzDomain.PermissionOrderRead}, // Staff has permission for shopID, but order items are in unassignedShopID
		},
	}
	ctx := authzSvc.WithActor(req.Context(), actor)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orderID", orderID.String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	err := handler.GetOrderTrackingForStaff(w, req)

	if err == nil {
		t.Fatalf("expected forbidden error for staff missing shop order:read permission, got nil")
	}
	appErr, ok := err.(*apperrors.AppError)
	if !ok || appErr.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 Forbidden, got %v", err)
	}
}

func TestUpdateOrderStatus_MultiShop_EnforcesAllShopRules(t *testing.T) {
	staffID := uuid.New()
	shop1ID := uuid.New()
	shop2ID := uuid.New()
	orderID := uuid.New()

	order := &orderDomain.Order{
		ID:    orderID,
		Total: 100000,
	}
	items := []orderDomain.OrderItem{
		{ID: uuid.New(), OrderID: orderID, ShopID: shop1ID},
		{ID: uuid.New(), OrderID: orderID, ShopID: shop2ID},
	}

	orderRepoMock := &hMockOrderRepo{order: order}
	orderItemRepoMock := &hMockOrderItemRepo{items: items}

	getOrderUC := usecase.NewGetOrderUsecase(
		&hMockExecutor{},
		orderRepoMock,
		orderItemRepoMock,
		&hMockPaymentRepo{},
		&hMockPaymentChannelDataRepo{},
		&hMockShipmentRepo{},
		&hMockShipmentEventRepo{},
	)

	handler := &orderHandler{
		getOrder: getOrderUC,
	}

	actor := &authzDomain.Actor{
		Type:    authenDomain.AccountTypeStaff,
		StaffID: &staffID,
		Permissions: map[uuid.UUID][]string{
			shop1ID: {authzDomain.PermissionOrderUpdateStatus},
			shop2ID: {authzDomain.PermissionOrderUpdateStatus},
		},
		Rules: map[uuid.UUID]map[string]any{
			shop1ID: {"allowed_statuses": []string{"confirmed", "shipped"}},
			shop2ID: {"allowed_statuses": []string{"confirmed"}}, // Shop 2 forbids 'shipped'
		},
	}

	body, _ := json.Marshal(map[string]any{
		"status": "shipped",
	})
	req := httptest.NewRequest(http.MethodPatch, "/orders/"+orderID.String()+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := authzSvc.WithActor(req.Context(), actor)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orderID", orderID.String())
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	err := handler.UpdateOrderStatus(w, req)

	if err == nil {
		t.Fatalf("expected forbidden error when transition is forbidden by shop 2 rule, got nil")
	}
	appErr, ok := err.(*apperrors.AppError)
	if !ok || appErr.StatusCode != http.StatusForbidden {
		t.Errorf("expected status 403 Forbidden, got %v", err)
	}
}
