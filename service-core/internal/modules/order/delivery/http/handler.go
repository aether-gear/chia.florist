package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	"service-core/internal/modules/order/usecase"

	"github.com/google/uuid"
)

type orderHandler struct {
	findOrders  *usecase.FindOrdersUsecase
	createOrder *usecase.CreateOrderUsecase
}

func NewOrderHandler(
	findOrders *usecase.FindOrdersUsecase,
	createOrder *usecase.CreateOrderUsecase,
) *orderHandler {
	return &orderHandler{
		findOrders:  findOrders,
		createOrder: createOrder,
	}
}

func (h *orderHandler) FindOrders(w http.ResponseWriter, r *http.Request) error {
	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}

	if actor.Type != authenDomain.AccountTypeMerchant {
		return apperrors.NewForbidden("forbidden: merchant account required")
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
	userIDStr := apphttp.Query(r, "user_id")
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

	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid user id")
		}
		input.UserID = &userID
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

		results[i] = orderResponse{
			ID:          o.Order.ID.String(),
			Number:      o.Order.Number,
			UserID:      o.Order.UserID.String(),
			AddressID:   o.Order.AddressID.String(),
			Status:      string(o.Order.Status),
			Subtotal:    o.Order.Subtotal,
			ShippingFee: o.Order.ShippingFee,
			Total:       o.Order.Total,
			CreatedAt:   o.Order.CreatedAt,
			UpdatedAt:   o.Order.UpdatedAt,
			Items:       items,
		}
	}

	response := map[string]any{
		"orders": results,
		"page":   page,
		"limit":  limit,
		"total":  total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	userID := authCtx.UserID

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
		AddressID:       parsedAddressID,
		PaymentMethodID: parsedPaymentMethodID,
		IsManual:        req.SelectedPayment.IsManual,
		Shops:           shopsInput,
	}

	result, err :=
		h.createOrder.Execute(r.Context(), input)
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
