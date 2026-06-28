package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	"service-core/internal/modules/order/usecase"
	paymentDomain "service-core/internal/modules/payment/domain"
	shipmentDomain "service-core/internal/modules/shipment/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type orderHandler struct {
	findOrders  *usecase.FindOrdersUsecase
	getOrder    *usecase.GetOrderUsecase
	createOrder *usecase.CreateOrderUsecase
}

func NewOrderHandler(
	findOrders *usecase.FindOrdersUsecase,
	getOrder *usecase.GetOrderUsecase,
	createOrder *usecase.CreateOrderUsecase,
) *orderHandler {
	return &orderHandler{
		findOrders:  findOrders,
		getOrder:    getOrder,
		createOrder: createOrder,
	}
}

// buildOrderResponse converts an OrderSearchResult into the wire response.
func buildOrderResponse(o usecase.OrderSearchResult) orderResponse {
	items := make([]orderItemResponse, len(o.Items))
	for j, item := range o.Items {
		items[j] = orderItemResponse{
			ID:               item.ID.String(),
			ProductID:        item.ProductID.String(),
			ProductName:      item.ProductName,
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			Subtotal:         item.Subtotal,
			ShopID:           item.ShopID.String(),
			ShopName:         item.ShopName,
			CourierCode:      item.CourierCode,
			CourierService:   item.CourierService,
			ShippingFeeTotal: item.ShippingFee,
		}
	}

	resp := orderResponse{
		ID:          o.Order.ID.String(),
		Number:      o.Order.Number,
		CustomerID:  o.Order.CustomerID.String(),
		AddressID:   o.Order.AddressID.String(),
		Status:      string(o.Order.Status),
		Subtotal:    o.Order.Subtotal,
		ShippingFee: o.Order.ShippingFee,
		Total:       o.Order.Total,
		CreatedAt:   o.Order.CreatedAt,
		UpdatedAt:   o.Order.UpdatedAt,
		Items:       items,
	}

	if o.Payment != nil {
		resp.Payment = mapPaymentDetail(o.Payment)
	}
	if o.Shipment != nil {
		resp.Shipment = mapShipmentDetail(o.Shipment)
	}

	return resp
}

func mapPaymentDetail(p *paymentDomain.Payment) *paymentDetailResponse {
	return &paymentDetailResponse{
		ID:        p.ID.String(),
		Status:    string(p.Status),
		Provider:  p.Provider,
		Amount:    p.Amount,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
}

func mapShipmentDetail(s *shipmentDomain.Shipment) *shipmentDetailResponse {
	return &shipmentDetailResponse{
		ID:             s.ID.String(),
		Status:         string(s.Status),
		Courier:        s.Courier,
		Service:        s.Service,
		TrackingNumber: s.TrackingNumber,
		Cost:           s.Cost,
		CreatedAt:      s.CreatedAt,
	}
}

func (h *orderHandler) FindOrders(w http.ResponseWriter, r *http.Request) error {
	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}

	if actor.Type != authenDomain.AccountTypeStaff {
		return apperrors.NewForbidden("forbidden: staff account required")
	}

	page := apphttp.QueryIntDefault(r, "page", 1)
	if page <= 0 {
		page = 1
	}
	limit := apphttp.QueryIntDefault(r, "limit", 10)
	if limit <= 0 {
		limit = 10
	}

	sort := apphttp.Query(r, "sort")
	idStr := apphttp.Query(r, "id")
	number := apphttp.Query(r, "number")
	customerIDStr := apphttp.Query(r, "customer_id")
	status := apphttp.Query(r, "status")

	input := usecase.FindOrdersInput{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}

	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid order id")
		}
		input.ID = &id
	}

	if number != "" {
		input.Number = &number
	}

	if customerIDStr != "" {
		customerID, err := uuid.Parse(customerIDStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid customer id")
		}
		input.CustomerID = &customerID
	}

	if status != "" {
		input.Status = &status
	}

	orders, total, err := h.findOrders.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	results := make([]orderResponse, len(orders))
	for i, o := range orders {
		results[i] = buildOrderResponse(o)
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]any{
		"orders": results,
		"page":   page,
		"limit":  limit,
		"total":  total,
	})
	return nil
}

// GetOrder handles GET /orders/{orderID} — staff-only, returns a single order with full detail.
func (h *orderHandler) GetOrder(w http.ResponseWriter, r *http.Request) error {
	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	if actor.Type != authenDomain.AccountTypeStaff {
		return apperrors.NewForbidden("forbidden: staff account required")
	}

	orderIDStr := chi.URLParam(r, "orderID")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid order id")
	}

	result, err := h.getOrder.Execute(r.Context(), usecase.GetOrderInput{
		OrderID: orderID,
	})
	if err != nil {
		return err
	}
	if result == nil {
		return apperrors.NewNotFound("order not found")
	}

	apphttp.WriteJSON(w, http.StatusOK, buildOrderResponse(usecase.OrderSearchResult{
		Order:    result.Order,
		Items:    result.Items,
		Payment:  result.Payment,
		Shipment: result.Shipment,
	}))
	return nil
}

// ListMyOrders handles GET /users/me/orders — customer-only, returns the caller's orders with detail.
func (h *orderHandler) ListMyOrders(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	customerID := *authCtx.CustomerID

	page := apphttp.QueryIntDefault(r, "page", 1)
	if page <= 0 {
		page = 1
	}
	limit := apphttp.QueryIntDefault(r, "limit", 10)
	if limit <= 0 {
		limit = 10
	}

	sort := apphttp.Query(r, "sort")
	status := apphttp.Query(r, "status")

	input := usecase.FindOrdersInput{
		Page:       page,
		Limit:      limit,
		Sort:       sort,
		CustomerID: &customerID,
	}

	if status != "" {
		input.Status = &status
	}

	orders, total, err := h.findOrders.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	results := make([]orderResponse, len(orders))
	for i, o := range orders {
		results[i] = buildOrderResponse(o)
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]any{
		"orders": results,
		"page":   page,
		"limit":  limit,
		"total":  total,
	})
	return nil
}

// GetMyOrder handles GET /users/me/orders/{orderID} — customer-only, returns their own order with full detail.
func (h *orderHandler) GetMyOrder(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	customerID := *authCtx.CustomerID

	orderIDStr := chi.URLParam(r, "orderID")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid order id")
	}

	result, err := h.getOrder.Execute(r.Context(), usecase.GetOrderInput{
		OrderID:    orderID,
		CustomerID: &customerID,
	})
	if err != nil {
		return err
	}
	if result == nil {
		return apperrors.NewNotFound("order not found")
	}

	apphttp.WriteJSON(w, http.StatusOK, buildOrderResponse(usecase.OrderSearchResult{
		Order:    result.Order,
		Items:    result.Items,
		Payment:  result.Payment,
		Shipment: result.Shipment,
	}))
	return nil
}

// CreateOrder handles POST /order — customer-only.
func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	userID := authCtx.UserID
	customerID := *authCtx.CustomerID

	var req createOrderRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	parsedAddressID, err := uuid.Parse(req.AddressID)
	if err != nil {
		return apperrors.NewBadRequest("invalid address id")
	}

	parsedPaymentMethodID, err := uuid.Parse(req.SelectedPayment.ID)
	if err != nil {
		return apperrors.NewBadRequest("invalid payment method id")
	}

	var shopsInput []usecase.OrderShopInput
	for _, shopReq := range req.Shops {
		parsedShopID, err := uuid.Parse(shopReq.ShopID)
		if err != nil {
			return apperrors.NewBadRequest("invalid shop id")
		}
		if shopReq.ShopName == "" {
			return apperrors.NewBadRequest("invalid shop name")
		}

		var itemsInput []usecase.OrderItemInput
		for _, itemReq := range shopReq.Items {
			if itemReq.ProductName == "" {
				return apperrors.NewBadRequest("invalid product name")
			}
			if itemReq.Quantity <= 0 {
				return apperrors.NewBadRequest("invalid quantity")
			}

			parsedProductID, err := uuid.Parse(itemReq.ProductID)
			if err != nil {
				return apperrors.NewBadRequest("invalid product id")
			}

			itemsInput = append(
				itemsInput,
				usecase.OrderItemInput{
					ProductID:   parsedProductID,
					ProductName: itemReq.ProductName,
					Quantity:    itemReq.Quantity,
				},
			)
		}

		var courierInput *usecase.OrderCourierInput
		if shopReq.Courier != nil {
			courierInput = &usecase.OrderCourierInput{
				Code:    shopReq.Courier.Code,
				Service: shopReq.Courier.Service,
			}
		}

		shopsInput = append(
			shopsInput,
			usecase.OrderShopInput{
				ShopID:   parsedShopID,
				ShopName: shopReq.ShopName,
				Courier:  courierInput,
				Items:    itemsInput,
			},
		)
	}

	input := usecase.CreateOrderInput{
		UserID:          userID,
		CustomerID:      customerID,
		AddressID:       parsedAddressID,
		PaymentMethodID: parsedPaymentMethodID,
		IsManual:        req.SelectedPayment.IsManual,
		Shops:           shopsInput,
	}

	result, err := h.createOrder.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	resp := createOrderResponse{
		OrderID:     result.OrderID.String(),
		Instruction: result.Instruction,
	}

	if result.PaymentAccount != nil {
		resp.PaymentAccount = &createOrderPaymentAccountResponse{
			AccountName:   result.PaymentAccount.AccountName,
			AccountNumber: result.PaymentAccount.AccountNumber,
			PhoneNumber:   result.PaymentAccount.PhoneNumber,
			QRString:      result.PaymentAccount.QRString,
		}
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}
