package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authenDomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/order/usecase"

	"github.com/google/uuid"
)

type orderHandler struct {
	createOrder *usecase.CreateOrderUsecase
}

func NewOrderHandler(
	createOrder *usecase.CreateOrderUsecase,
) *orderHandler {
	return &orderHandler{
		createOrder: createOrder,
	}
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
