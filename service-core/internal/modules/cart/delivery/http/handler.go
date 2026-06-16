package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authdomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/cart/usecase"

	"github.com/google/uuid"
)

type CartHandler struct {
	addItem    *usecase.AddItemUsecase
	getCart    *usecase.GetCartUsecase
	updateItem *usecase.UpdateItemUsecase
	removeItem *usecase.RemoveItemUsecase
	checkout   *usecase.CheckoutUsecase
}

func NewCartHandler(
	addItem *usecase.AddItemUsecase,
	getCart *usecase.GetCartUsecase,
	updateItem *usecase.UpdateItemUsecase,
	removeItem *usecase.RemoveItemUsecase,
	checkout *usecase.CheckoutUsecase,
) *CartHandler {
	return &CartHandler{
		addItem:    addItem,
		getCart:    getCart,
		updateItem: updateItem,
		removeItem: removeItem,
		checkout:   checkout,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	result, err := h.getCart.Execute(r.Context(), authCtx.UserID)
	if err != nil {
		return err
	}
	if result == nil || result.Cart == nil {
		return apperrors.NewNotFound("cart not found")
	}

	var total int64
	items := make([]cartItemView, 0, len(result.Cart.Items))

	for _, item := range result.Cart.Items {
		if result.Products == nil {
			continue
		}

		productData, ok := result.Products[item.ProductID]
		if !ok {
			continue
		}

		price := productData.Product.Price
		quantity := item.Quantity
		subtotal := price * int64(quantity)

		image := productImageResponse{}
		if productData.Images.Thumbnail != "" {
			image.Thumbnail = &productData.Images.Thumbnail
		}

		items = append(items, cartItemView{
			ProductID: item.ProductID,
			ShopID:    item.ShopID,
			Name:      productData.Product.Name,
			Price:     price,
			Quantity:  quantity,
			Subtotal:  subtotal,
			Image:     image,
		})

		total += subtotal
	}

	response := cartResponse{
		CartID: result.Cart.ID,
		Items:  items,
		Total:  total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) error {
	var req addItemRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	if req.Quantity <= 0 {
		return apperrors.NewBadRequest("invalid quantity")
	}

	input := usecase.AddItemInput{
		UserID:    authCtx.UserID,
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  req.Quantity,
	}

	if err := h.addItem.Execute(r.Context(), input); err != nil {
		return err
	}

	response := map[string]string{
		"message": "item added",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) error {
	var req updateItemRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	productID, err := apphttp.ParamUUID(r, "productID")
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	shopID, err := apphttp.ParamUUID(r, "shopID")
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	if req.Quantity < 0 {
		return apperrors.NewBadRequest("invalid quantity")
	}

	input := usecase.UpdateItemInput{
		UserID:    authCtx.UserID,
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  req.Quantity,
	}

	if err := h.updateItem.Execute(r.Context(), input); err != nil {
		return err
	}

	response := map[string]string{
		"message": "item updated",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	productID, err := apphttp.ParamUUID(r, "productID")
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	shopID, err := apphttp.ParamUUID(r, "shopID")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	input := usecase.RemoveItemInput{
		UserID:    authCtx.UserID,
		ProductID: productID,
		ShopID:    shopID,
	}

	if err := h.removeItem.Execute(r.Context(), input); err != nil {
		return err
	}

	response := map[string]string{
		"message": "item removed",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) error {
	authCtx, input, err := h.parseCheckoutInput(r, false)
	if err != nil {
		return err
	}

	result, err := h.checkout.Execute(
		r.Context(),
		*authCtx,
		input,
	)
	if err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, toCheckoutResponse(*result))
	return nil
}

func (h *CartHandler) CheckoutEstimate(w http.ResponseWriter, r *http.Request) error {
	authCtx, input, err := h.parseCheckoutInput(r, true)
	if err != nil {
		return err
	}

	result, err := h.checkout.Execute(
		r.Context(),
		*authCtx,
		input,
	)
	if err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, toCheckoutResponse(*result))
	return nil
}

func (h *CartHandler) parseCheckoutInput(
	r *http.Request,
	requireCourier bool,
) (*authdomain.AuthContext, usecase.CheckoutInput, error) {
	var req checkoutRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return nil, usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid request body")
	}

	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return nil, usecase.CheckoutInput{}, apperrors.NewUnauthorized("authentication required")
	}

	var addressID *uuid.UUID
	if req.AddressID != nil {
		parsed, err := uuid.Parse(*req.AddressID)
		if err != nil {
			return nil, usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid address id")
		}
		addressID = &parsed
	}

	var shopInput []usecase.CheckoutShopInput
	for _, shopReq := range req.Shops {
		shopID, err := uuid.Parse(shopReq.ShopID)
		if err != nil {
			return nil, usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid shop id")
		}

		var items []usecase.CheckoutItemInput
		for _, itemReq := range shopReq.Items {
			productID, err := uuid.Parse(itemReq.ProductID)
			if err != nil {
				return nil, usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid product id")
			}

			if itemReq.Quantity <= 0 {
				return nil, usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid quantity")
			}

			items = append(items, usecase.CheckoutItemInput{
				ProductID: productID,
				Quantity:  itemReq.Quantity,
			})
		}

		var courier *usecase.SelectedCourierInput
		if shopReq.Courier != nil {
			courier = &usecase.SelectedCourierInput{
				Code:    shopReq.Courier.Code,
				Service: shopReq.Courier.Service,
			}
		}

		if requireCourier && courier == nil {
			return nil, usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid courier")
		}

		shopInput = append(shopInput, usecase.CheckoutShopInput{
			ShopID:  shopID,
			Items:   items,
			Courier: courier,
		})
	}

	return authCtx, usecase.CheckoutInput{
		AddressID: addressID,
		ShopInput: shopInput,
	}, nil
}

func toCheckoutResponse(result usecase.CheckoutResult) checkoutResponse {
	var shopsResponse []shopResponse
	for _, shop := range result.Shops {
		var itemsResponse []checkoutItemResponse
		for _, item := range shop.Items {
			itemsResponse = append(itemsResponse, checkoutItemResponse{
				ProductID: item.ProductID,
				ShopID:    item.ShopID,
				Name:      item.Name,
				Price:     item.Price,
				Quantity:  item.Quantity,
				Subtotal:  item.Subtotal,
			})
		}

		var shippingResponse []checkoutCouriersResponse
		for _, courier := range shop.CostCouriers {
			shippingResponse = append(shippingResponse, checkoutCouriersResponse{
				Code:    courier.Code,
				Name:    courier.Name,
				Service: courier.Service,
				ETD:     courier.ETD,
				Fee:     courier.Fee,
			})
		}

		shopReponse := shopResponse{
			ShopID:       shop.ShopID,
			Subtotal:     shop.Subtotal,
			Total:        &shop.Total,
			Items:        itemsResponse,
			CostCouriers: shippingResponse,
		}
		if shop.SelectedCourier != nil {
			shopReponse.SelectedCourier = &selectedCourierResponse{
				Code:    shop.SelectedCourier.Code,
				Service: shop.SelectedCourier.Service,
				Fee:     shop.SelectedCourier.Fee,
			}
		}
		shopsResponse = append(shopsResponse, shopReponse)
	}

	resp := checkoutResponse{
		Address: checkoutAddressResponse{
			ID:            result.Address.ID,
			RecipientName: result.Address.RecipientName,
			Phone:         result.Address.Phone,
			FullAddress:   result.Address.FullAddress,
		},
		Shops:         shopsResponse,
		TotalShipping: result.TotalShippingFee,
		Subtotal:      result.Subtotal,
	}

	if result.TotalAll > 0 {
		resp.TotalAll = &result.TotalAll
	}

	return resp
}
