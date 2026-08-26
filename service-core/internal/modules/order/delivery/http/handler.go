package http

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzDomain "service-core/internal/modules/authorization/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	cartDomain "service-core/internal/modules/cart/domain"
	orderDomain "service-core/internal/modules/order/domain"
	"service-core/internal/modules/order/usecase"
	paymentDomain "service-core/internal/modules/payment/domain"
	shipmentDomain "service-core/internal/modules/shipment/domain"
	shopUsecase "service-core/internal/modules/shop/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type orderHandler struct {
	findOrders           *usecase.FindOrdersUsecase
	getOrder             *usecase.GetOrderUsecase
	createOrder          *usecase.CreateOrderUsecase
	updateOrderStatus    *usecase.UpdateOrderStatusUsecase
	dispatchShopShipment *usecase.DispatchShopShipmentUsecase
	getOrderTracking     *usecase.GetOrderTrackingUsecase
	getShop              *shopUsecase.GetShopUsecase
}

func NewOrderHandler(
	findOrders *usecase.FindOrdersUsecase,
	getOrder *usecase.GetOrderUsecase,
	createOrder *usecase.CreateOrderUsecase,
	updateOrderStatus *usecase.UpdateOrderStatusUsecase,
	dispatchShopShipment *usecase.DispatchShopShipmentUsecase,
	getOrderTracking *usecase.GetOrderTrackingUsecase,
	getShop *shopUsecase.GetShopUsecase,
) *orderHandler {
	return &orderHandler{
		findOrders:           findOrders,
		getOrder:             getOrder,
		createOrder:          createOrder,
		updateOrderStatus:    updateOrderStatus,
		dispatchShopShipment: dispatchShopShipment,
		getOrderTracking:     getOrderTracking,
		getShop:              getShop,
	}
}

// buildOrderResponse converts an OrderSearchResult into the wire response.
func buildOrderResponse(o usecase.OrderSearchResult) orderResponse {
	items := make([]orderItemResponse, len(o.Items))
	for j, item := range o.Items {
		var productIDStr *string
		if item.ProductID != nil {
			s := item.ProductID.String()
			productIDStr = &s
		}
		var shipmentIDStr *string
		if item.ShipmentID != nil {
			s := item.ShipmentID.String()
			shipmentIDStr = &s
		}
		variantType := string(item.ProductVariantType)
		if variantType == "" {
			if item.ProductID == nil {
				variantType = "custom"
			} else {
				variantType = "standard"
			}
		}

		var customDesignResp *orderItemCustomDesignResponse
		if o.CustomDesigns != nil {
			if cd, ok := o.CustomDesigns[item.ID]; ok {
				customDesignResp = &orderItemCustomDesignResponse{
					Version:         cd.Version,
					PhysicalSizeID:  cd.PhysicalSizeID,
					PreviewURL:      cd.PreviewURL,
					HeaderTextUpper: cd.HeaderTextUpper,
					BodyTextUpper:   cd.BodyTextUpper,
					HeaderTextLower: cd.HeaderTextLower,
					BodyTextLower:   cd.BodyTextLower,
					DesignSnapshot:  cd.DesignSnapshot,
				}
			}
		}

		var itemOptsResp *orderItemOptionsResponse
		if item.ProductID != nil && variantType == "standard" {
			normOpts := item.ItemOptions.Normalized()
			itemOptsResp = &orderItemOptionsResponse{
				Size:   normOpts.Size,
				Jambul: normOpts.Jambul,
			}
		}

		items[j] = orderItemResponse{
			ID:                 item.ID.String(),
			ShipmentID:         shipmentIDStr,
			ProductID:          productIDStr,
			ProductVariantType: variantType,
			IsCustom:           item.ProductID == nil || variantType == "custom",
			ProductName:        item.ProductName,
			Quantity:           item.Quantity,
			UnitPrice:          item.UnitPrice,
			Subtotal:           item.Subtotal,
			ShopID:             item.ShopID.String(),
			ShopName:           item.ShopName,
			CourierCode:        item.CourierCode,
			CourierService:     item.CourierService,
			ShippingFeeTotal:   item.ShippingFee,
			ItemOptions:        itemOptsResp,
			CustomDesign:       customDesignResp,
		}
	}

	shipments := make([]shipmentDetailResponse, len(o.Shipments))
	for i := range o.Shipments {
		shipments[i] = *mapShipmentDetail(&o.Shipments[i])
	}

	resp := orderResponse{
		ID:                o.Order.ID.String(),
		Number:            o.Order.Number,
		CustomerID:        o.Order.CustomerID.String(),
		AddressID:         o.Order.AddressID.String(),
		Status:            string(o.Order.Status),
		Subtotal:          o.Order.Subtotal,
		ShippingFee:       o.Order.ShippingFee,
		Total:             o.Order.Total,
		ConfirmedAt:       o.Order.ConfirmedAt,
		HandlingExpiresAt: o.Order.HandlingExpiresAt,
		CreatedAt:         o.Order.CreatedAt,
		UpdatedAt:         o.Order.UpdatedAt,
		Items:             items,
		Shipments:         shipments,
	}

	if o.Payment != nil {
		resp.Payment = mapPaymentDetail(o.Payment, o.ChannelData)
	}
	if o.Shipment != nil {
		resp.Shipment = mapShipmentDetail(o.Shipment)
	} else if len(shipments) > 0 {
		resp.Shipment = &shipments[0]
	}
	if o.Address != nil {
		resp.Address = &orderAddressResponse{
			ID:           o.Address.ID.String(),
			CustomerID:   o.Address.CustomerID.String(),
			ReceiverName: o.Address.ReceiverName,
			Phone:        o.Address.Phone,
			IsDefault:    o.Address.IsDefault,
			ProvinceID:   o.Address.Detail.ProvinceID,
			CityID:       o.Address.Detail.CityID,
			DistrictID:   o.Address.Detail.DistrictID,
			VillageID:    o.Address.Detail.VillageID,
			FullAddress:  o.Address.Detail.FullAddress,
			PostalCode:   o.Address.Detail.PostalCode,
		}
	}

	return resp
}

func mapPaymentDetail(p *paymentDomain.Payment, cd *paymentDomain.PaymentChannelData) *paymentDetailResponse {
	resp := &paymentDetailResponse{
		ID:        p.ID.String(),
		Status:    string(p.Status),
		Provider:  p.Provider,
		Amount:    p.Amount,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
	if cd != nil {
		resp.ChannelData = &paymentChannelDataResponse{
			ChannelType: string(cd.ChannelType),
			DisplayName: cd.DisplayName,
			ActionURL:   cd.ActionURL,
			ExpiresAt:   cd.ExpiresAt,
		}
	}

	return resp
}

func mapShipmentDetail(s *shipmentDomain.Shipment) *shipmentDetailResponse {
	events := make([]shipmentEventResponse, len(s.Events))
	for i, e := range s.Events {
		events[i] = shipmentEventResponse{
			ID:          e.ID.String(),
			Status:      e.Status,
			Description: e.Description,
			Location:    e.Location,
			Timestamp:   e.Timestamp,
		}
	}

	itemIDs := make([]string, len(s.ItemIDs))
	for i, id := range s.ItemIDs {
		itemIDs[i] = id.String()
	}

	return &shipmentDetailResponse{
		ID:                s.ID.String(),
		OrderID:           s.OrderID.String(),
		Status:            string(s.Status),
		FulfillmentMethod: string(s.FulfillmentMethod),
		Courier:           s.Courier,
		Service:           s.Service,
		TrackingNumber:    s.TrackingNumber,
		Cost:              s.Cost,
		CreatedAt:         s.CreatedAt,
		Events:            events,
		ItemIDs:           itemIDs,
	}
}

func (h *orderHandler) resolveShopFilter(r *http.Request) (*uuid.UUID, bool, error) {
	shopIDStr := apphttp.Query(r, "shop_id")
	shopSlug := apphttp.Query(r, "shop_slug")
	shopParam := apphttp.Query(r, "shop")

	if shopIDStr == "all" || shopSlug == "all" || shopParam == "all" {
		return nil, false, nil
	}

	targetIDStr := ""
	targetSlug := ""

	if shopIDStr != "" {
		targetIDStr = shopIDStr
	} else if shopSlug != "" {
		targetSlug = shopSlug
	} else if shopParam != "" {
		if parsed, err := uuid.Parse(shopParam); err == nil {
			return &parsed, true, nil
		}
		targetSlug = shopParam
	}

	if targetIDStr != "" {
		id, err := uuid.Parse(targetIDStr)
		if err != nil {
			return nil, true, apperrors.NewBadRequest("invalid shop id")
		}
		return &id, true, nil
	}

	if targetSlug != "" {
		if targetSlug == "all" {
			return nil, false, nil
		}
		if h.getShop == nil {
			return nil, true, apperrors.NewInternal(errors.New("shop filter service unavailable"))
		}
		shop, err := h.getShop.GetBySlug(r.Context(), targetSlug)
		if err != nil {
			return nil, true, err
		}
		if shop == nil {
			return nil, true, nil
		}
		return &shop.ID, true, nil
	}

	return nil, false, nil
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
	} else if statusesParam := apphttp.Query(r, "statuses"); statusesParam != "" {
		var parsedStatuses []string
		parts := strings.Split(statusesParam, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				parsedStatuses = append(parsedStatuses, trimmed)
			}
		}
		input.Statuses = parsedStatuses
	}

	fromDateStr := apphttp.Query(r, "from_date")
	if fromDateStr != "" {
		if t, err := time.Parse(time.RFC3339, fromDateStr); err == nil {
			input.FromDate = &t
		} else if t, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			input.FromDate = &t
		} else {
			return apperrors.NewBadRequest("invalid from_date format")
		}
	}

	toDateStr := apphttp.Query(r, "to_date")
	if toDateStr != "" {
		if t, err := time.Parse(time.RFC3339, toDateStr); err == nil {
			input.ToDate = &t
		} else if t, err := time.Parse("2006-01-02", toDateStr); err == nil {
			endOfDay := t.Add(24*time.Hour - time.Nanosecond)
			input.ToDate = &endOfDay
		} else {
			return apperrors.NewBadRequest("invalid to_date format")
		}
	}

	shopID, shopSpecified, err := h.resolveShopFilter(r)
	if err != nil {
		return err
	}
	if shopSpecified {
		if shopID == nil {
			apphttp.WriteJSON(w, http.StatusOK, map[string]any{
				"orders": []orderResponse{},
				"page":   page,
				"limit":  limit,
				"total":  0,
			})
			return nil
		}
		input.ShopID = shopID
	}

	if actor.StaffID != nil && !actor.IsSuperAdmin() {
		allAssigned := actor.GetAssignedShopIDs()
		var assignedIDs []uuid.UUID
		for _, sID := range allAssigned {
			if actor.HasPermission(sID, authzDomain.PermissionOrderRead) {
				assignedIDs = append(assignedIDs, sID)
			}
		}

		if len(assignedIDs) == 0 {
			apphttp.WriteJSON(w, http.StatusOK, map[string]any{
				"orders": []orderResponse{},
				"page":   page,
				"limit":  limit,
				"total":  0,
			})
			return nil
		}

		if input.ShopID != nil {
			if !actor.HasPermission(*input.ShopID, authzDomain.PermissionOrderRead) {
				return apperrors.NewForbidden("forbidden: missing order:read permission for this shop")
			}
		} else {
			input.ShopIDs = assignedIDs
		}
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

	if actor.StaffID != nil && !actor.IsSuperAdmin() {
		var permittedItems []orderDomain.OrderItem
		for _, item := range result.Items {
			if actor.HasPermission(item.ShopID, authzDomain.PermissionOrderRead) {
				permittedItems = append(permittedItems, item)
			}
		}
		if len(permittedItems) == 0 {
			return apperrors.NewForbidden("forbidden: missing order:read permission for this shop's order")
		}
		result.Items = permittedItems

		var permittedShipments []shipmentDomain.Shipment
		for _, s := range result.Shipments {
			hasItem := false
			for _, itm := range permittedItems {
				if itm.ShipmentID != nil && *itm.ShipmentID == s.ID {
					hasItem = true
					break
				}
			}
			if hasItem {
				permittedShipments = append(permittedShipments, s)
			}
		}
		result.Shipments = permittedShipments
		if len(permittedShipments) > 0 {
			result.Shipment = &permittedShipments[0]
		} else {
			result.Shipment = nil
		}
	}

	resp := buildOrderResponse(usecase.OrderSearchResult{
		Order:         result.Order,
		Items:         result.Items,
		CustomDesigns: result.CustomDesigns,
		Payment:       result.Payment,
		ChannelData:   result.ChannelData,
		Shipment:      result.Shipment,
		Shipments:     result.Shipments,
	})

	apphttp.WriteJSON(w, http.StatusOK, resp)
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
	} else if statusesParam := apphttp.Query(r, "statuses"); statusesParam != "" {
		var parsedStatuses []string
		parts := strings.Split(statusesParam, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				parsedStatuses = append(parsedStatuses, trimmed)
			}
		}
		input.Statuses = parsedStatuses
	}

	fromDateStr := apphttp.Query(r, "from_date")
	if fromDateStr != "" {
		if t, err := time.Parse(time.RFC3339, fromDateStr); err == nil {
			input.FromDate = &t
		} else if t, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			input.FromDate = &t
		} else {
			return apperrors.NewBadRequest("invalid from_date format")
		}
	}

	toDateStr := apphttp.Query(r, "to_date")
	if toDateStr != "" {
		if t, err := time.Parse(time.RFC3339, toDateStr); err == nil {
			input.ToDate = &t
		} else if t, err := time.Parse("2006-01-02", toDateStr); err == nil {
			endOfDay := t.Add(24*time.Hour - time.Nanosecond)
			input.ToDate = &endOfDay
		} else {
			return apperrors.NewBadRequest("invalid to_date format")
		}
	}

	shopID, shopSpecified, err := h.resolveShopFilter(r)
	if err != nil {
		return err
	}
	if shopSpecified {
		if shopID == nil {
			apphttp.WriteJSON(w, http.StatusOK, map[string]any{
				"orders": []orderResponse{},
				"page":   page,
				"limit":  limit,
				"total":  0,
			})
			return nil
		}
		input.ShopID = shopID
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

	resp := buildOrderResponse(usecase.OrderSearchResult{
		Order:         result.Order,
		Items:         result.Items,
		CustomDesigns: result.CustomDesigns,
		Payment:       result.Payment,
		ChannelData:   result.ChannelData,
		Shipment:      result.Shipment,
		Shipments:     result.Shipments,
	})

	apphttp.WriteJSON(w, http.StatusOK, resp)
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
			if itemReq.Quantity <= 0 {
				return apperrors.NewBadRequest("invalid quantity")
			}

			var productID *uuid.UUID
			var cartItemID *uuid.UUID

			isCustom := (itemReq.IsCustom != nil && *itemReq.IsCustom) ||
				itemReq.ProductVariantType == "custom" ||
				itemReq.ItemType == "custom" ||
				len(itemReq.CustomDesign) > 0

			if itemReq.ProductID != nil &&
				*itemReq.ProductID != "" &&
				*itemReq.ProductID != "null" &&
				*itemReq.ProductID != "undefined" &&
				*itemReq.ProductID != "custom" {

				if parsed, err := uuid.Parse(*itemReq.ProductID); err == nil {
					productID = &parsed
				} else if !isCustom {
					return apperrors.NewBadRequest("invalid product id")
				}
			}

			if itemReq.CartItemID != nil &&
				*itemReq.CartItemID != "" &&
				*itemReq.CartItemID != "null" &&
				*itemReq.CartItemID != "undefined" {

				if parsed, err := uuid.Parse(*itemReq.CartItemID); err == nil {
					cartItemID = &parsed
				}
			}

			if !isCustom && productID == nil && cartItemID == nil {
				return apperrors.NewBadRequest("invalid product id")
			}

			productName := itemReq.ProductName
			if productName == "" {
				if isCustom {
					productName = "(Custom Flower Board)"
				} else {
					return apperrors.NewBadRequest("invalid product name")
				}
			}

			var opt cartDomain.ItemOptions
			if itemReq.ItemOptions != nil {
				opt.Size = itemReq.ItemOptions.Size
				opt.Jambul = itemReq.ItemOptions.Jambul
			}
			if itemReq.Size != "" && opt.Size == "" {
				opt.Size = itemReq.Size
			}
			if itemReq.Jambul != "" && opt.Jambul == "" {
				opt.Jambul = itemReq.Jambul
			}

			itemsInput = append(
				itemsInput,
				usecase.OrderItemInput{
					ProductID:    productID,
					CartItemID:   cartItemID,
					IsCustom:     isCustom,
					ItemOptions:  opt.Normalized(),
					CustomDesign: itemReq.CustomDesign,
					ProductName:  productName,
					Quantity:     itemReq.Quantity,
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

	if result.ChannelData != nil {
		resp.ChannelData = &paymentChannelDataResponse{
			ChannelType: string(result.ChannelData.ChannelType),
			DisplayName: result.ChannelData.DisplayName,
			ActionURL:   result.ChannelData.ActionURL,
			ExpiresAt:   result.ChannelData.ExpiresAt,
		}
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

// UpdateOrderStatus handles PATCH /orders/{orderID}/status — staff-only.
// Transitions the order to the requested status. When the target status is
// "shipped", the configured LogisticsProvider creates a shipment record.
// In Komerce mode the provider calls the external API and returns a tracking
// number automatically. In manual mode the optional "tracking_number" field
// in the request body is used instead — no external call is made.
func (h *orderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) error {
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

	var req updateOrderStatusRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}
	if req.Status == "" {
		return apperrors.NewBadRequest("status is required")
	}

	if actor.StaffID != nil && !actor.IsSuperAdmin() {
		existingOrder, err := h.getOrder.Execute(r.Context(), usecase.GetOrderInput{
			OrderID: orderID,
		})
		if err != nil || existingOrder == nil {
			return apperrors.NewNotFound("order not found")
		}

		uniqueShops := make(map[uuid.UUID]bool)
		for _, item := range existingOrder.Items {
			if actor.HasPermission(item.ShopID, authzDomain.PermissionOrderUpdateStatus) {
				uniqueShops[item.ShopID] = true
			}
		}
		if len(uniqueShops) == 0 {
			return apperrors.NewForbidden("forbidden: missing order:update_status permission for this shop's order")
		}

		for shopID := range uniqueShops {
			if shopRules, exists := actor.Rules[shopID]; exists && shopRules != nil {
				if allowedRaw, ok := shopRules["allowed_statuses"]; ok {
					allowedList := parseStringSlice(allowedRaw)
					if len(allowedList) > 0 && !slices.Contains(allowedList, req.Status) {
						return apperrors.NewForbidden("forbidden: status transition to '" + req.Status + "' is not allowed by staff rule")
					}
				}

				if maxRaw, ok := shopRules["max_order_amount"]; ok {
					maxAmount := parseFloat64(maxRaw)
					if maxAmount > 0 && float64(existingOrder.Order.Total) > maxAmount {
						return apperrors.NewForbidden("forbidden: order total amount exceeds staff rule limit")
					}
				}
			}
		}
	}

	var shipmentsInput []usecase.ShipmentDispatchInput
	if len(req.Shipments) > 0 {
		for _, sReq := range req.Shipments {
			var itemUUIDs []uuid.UUID
			for _, idStr := range sReq.ItemIDs {
				parsed, err := uuid.Parse(idStr)
				if err != nil {
					return apperrors.NewBadRequest("invalid item id in shipments")
				}
				itemUUIDs = append(itemUUIDs, parsed)
			}
			shipmentsInput = append(shipmentsInput, usecase.ShipmentDispatchInput{
				FulfillmentMethod: sReq.FulfillmentMethod,
				Courier:           sReq.Courier,
				Service:           sReq.Service,
				TrackingNumber:    sReq.TrackingNumber,
				ItemIDs:           itemUUIDs,
			})
		}
	}

	result, err := h.updateOrderStatus.Execute(r.Context(), usecase.UpdateOrderStatusInput{
		OrderID:           orderID,
		Status:            orderDomain.OrderStatus(req.Status),
		TrackingNumber:    req.TrackingNumber,
		FulfillmentMethod: req.FulfillmentMethod,
		Shipments:         shipmentsInput,
	})
	if err != nil {
		return err
	}

	resp := buildOrderResponse(usecase.OrderSearchResult{
		Order:     result.Order,
		Items:     []orderDomain.OrderItem{},
		Shipment:  result.Shipment,
		Shipments: result.Shipments,
	})

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

// DispatchOrderShipment handles POST /orders/{orderID}/shipments — staff-only.
// Creates a shipment for a specific shop's order items.
func (h *orderHandler) DispatchOrderShipment(w http.ResponseWriter, r *http.Request) error {
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

	var req dispatchShopShipmentRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	if len(req.ItemIDs) == 0 {
		return apperrors.NewBadRequest("item_ids is required and must not be empty")
	}

	var itemUUIDs []uuid.UUID
	for _, idStr := range req.ItemIDs {
		parsed, err := uuid.Parse(idStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid item id in item_ids")
		}
		itemUUIDs = append(itemUUIDs, parsed)
	}

	if actor.StaffID != nil && !actor.IsSuperAdmin() {
		if !actor.HasPermission(shopID, authzDomain.PermissionOrderUpdateStatus) {
			return apperrors.NewForbidden("forbidden: missing order:update_status permission for this shop")
		}
	}

	res, err := h.dispatchShopShipment.Execute(r.Context(), usecase.DispatchShopShipmentInput{
		OrderID:           orderID,
		ShopID:            shopID,
		FulfillmentMethod: req.FulfillmentMethod,
		Courier:           req.Courier,
		Service:           req.Service,
		TrackingNumber:    req.TrackingNumber,
		ItemIDs:           itemUUIDs,
	})
	if err != nil {
		return err
	}

	resp := map[string]any{
		"order_id":          res.Order.ID.String(),
		"order_status":      string(res.Order.Status),
		"shipment_id":       res.Shipment.ID.String(),
		"all_items_shipped": res.AllItemsShipped,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *orderHandler) GetMyOrderTracking(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	customerID := *authCtx.CustomerID

	orderID, err := apphttp.ParamUUID(r, "orderID")
	if err != nil {
		return apperrors.NewBadRequest("invalid order id")
	}

	input := usecase.GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: customerID,
	}

	result, err := h.getOrderTracking.Execute(r.Context(), input)
	if err != nil {
		return err
	}
	if result == nil {
		return apperrors.NewNotFound("tracking information not found")
	}

	timeline := make([]trackingTimelineEventResponse, len(result.Timeline))
	for i, e := range result.Timeline {
		timeline[i] = trackingTimelineEventResponse{
			Status:      e.Status,
			Description: e.Description,
			Location:    e.Location,
			Timestamp:   e.Timestamp,
		}
	}

	if result.Warning != nil {
		w.Header().Set("X-Warning", *result.Warning)
	}

	resp := orderTrackingResponse{
		OrderID:        result.OrderID.String(),
		ShipmentID:     result.ShipmentID.String(),
		Courier:        result.Courier,
		TrackingNumber: result.TrackingNumber,
		Warning:        result.Warning,
		Timeline:       timeline,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

// GetOrderTrackingForStaff handles GET /orders/{orderID}/tracking — staff-only.
func (h *orderHandler) GetOrderTrackingForStaff(w http.ResponseWriter, r *http.Request) error {
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

	if actor.StaffID != nil && !actor.IsSuperAdmin() {
		existingOrder, err := h.getOrder.Execute(r.Context(), usecase.GetOrderInput{
			OrderID: orderID,
		})
		if err != nil || existingOrder == nil {
			return apperrors.NewNotFound("order not found")
		}

		hasAccess := false
		for _, item := range existingOrder.Items {
			if actor.HasPermission(item.ShopID, authzDomain.PermissionOrderRead) {
				hasAccess = true
				break
			}
		}
		if !hasAccess {
			return apperrors.NewForbidden("forbidden: missing order:read permission for this shop's order")
		}
	}

	input := usecase.GetOrderTrackingInput{
		OrderID:    orderID,
		CustomerID: uuid.Nil,
	}

	result, err := h.getOrderTracking.Execute(r.Context(), input)
	if err != nil {
		return err
	}
	if result == nil {
		return apperrors.NewNotFound("tracking information not found")
	}

	if result.Warning != nil {
		w.Header().Set("X-Warning", *result.Warning)
	}

	timeline := make([]trackingTimelineEventResponse, len(result.Timeline))
	for i, e := range result.Timeline {
		timeline[i] = trackingTimelineEventResponse{
			Status:      e.Status,
			Description: e.Description,
			Location:    e.Location,
			Timestamp:   e.Timestamp,
		}
	}

	resp := orderTrackingResponse{
		OrderID:        result.OrderID.String(),
		ShipmentID:     result.ShipmentID.String(),
		Courier:        result.Courier,
		TrackingNumber: result.TrackingNumber,
		Warning:        result.Warning,
		Timeline:       timeline,
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func parseStringSlice(val any) []string {
	var result []string
	switch v := val.(type) {
	case []string:
		return v
	case []any:
		for _, elem := range v {
			if s, ok := elem.(string); ok {
				result = append(result, s)
			}
		}
	}
	return result
}

func parseFloat64(val any) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case int:
		return float64(v)
	}
	return 0
}
