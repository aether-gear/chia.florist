package http

import (
	"encoding/json"
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authdomain "service-core/internal/modules/authentication/domain"
	domain "service-core/internal/modules/cart/domain"
	"service-core/internal/modules/cart/usecase"
	productDomain "service-core/internal/modules/product/domain"

	"github.com/google/uuid"
)

type CartHandler struct {
	addItem          *usecase.AddItemUsecase
	addCustomItem    *usecase.AddCustomItemUsecase
	getCart          *usecase.GetCartUsecase
	updateItem       *usecase.UpdateItemUsecase
	removeItem       *usecase.RemoveItemUsecase
	removeCustomItem *usecase.RemoveCustomItemUsecase
	changeItemShop   *usecase.ChangeItemShopUsecase
	checkout         *usecase.CheckoutUsecase
}

func NewCartHandler(
	addItem *usecase.AddItemUsecase,
	addCustomItem *usecase.AddCustomItemUsecase,
	getCart *usecase.GetCartUsecase,
	updateItem *usecase.UpdateItemUsecase,
	removeItem *usecase.RemoveItemUsecase,
	removeCustomItem *usecase.RemoveCustomItemUsecase,
	changeItemShop *usecase.ChangeItemShopUsecase,
	checkout *usecase.CheckoutUsecase,
) *CartHandler {
	return &CartHandler{
		addItem:          addItem,
		addCustomItem:    addCustomItem,
		getCart:          getCart,
		updateItem:       updateItem,
		removeItem:       removeItem,
		removeCustomItem: removeCustomItem,
		changeItemShop:   changeItemShop,
		checkout:         checkout,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	result, err := h.getCart.Execute(r.Context(), *authCtx.CustomerID)
	if err != nil {
		return err
	}
	if result == nil || result.Cart == nil {
		return apperrors.NewNotFound("cart not found")
	}

	var total int64
	items := make([]cartItemView, 0, len(result.Cart.Items))

	for _, item := range result.Cart.Items {
		var shopName, shopSlug string
		if result.Shops != nil {
			if s, ok := result.Shops[item.ShopID]; ok {
				shopName = s.Name
				shopSlug = s.Slug
			}
		}

		// Compute custom product price dynamically by backend-side
		if item.ProductVariantType == domain.ProductVariantTypeCustom {
			var price int64
			var subtotal int64
			if len(item.CustomDesign) > 0 {
				if parsedDesign, err := productDomain.ParseCustomDesignPayload(item.CustomDesign); err == nil {
					breakdown := productDomain.CalculateCustomProductPrice(
						*parsedDesign,
						productDomain.DefaultCustomPricingMatrix(),
					)
					price = breakdown.TotalPrice
					subtotal = price * int64(item.Quantity)
				}
			}

			imageResp := productImageResponse{}
			if len(item.CustomDesign) > 0 {
				var rawMap map[string]interface{}
				if err := json.Unmarshal(item.CustomDesign, &rawMap); err == nil {
					if assets, ok := rawMap["assets"].(map[string]interface{}); ok {
						if previewURL, ok := assets["previewUrl"].(string); ok && previewURL != "" {
							imageResp.Thumbnail = &previewURL
						}
					}
				}
			}

			items = append(items, cartItemView{
				CartItemID:         item.ID,
				ProductVariantType: string(item.ProductVariantType),
				ShopID:             item.ShopID,
				ShopName:           shopName,
				ShopSlug:           shopSlug,
				Name:               "(Custom Board)",
				Price:              price,
				Quantity:           item.Quantity,
				Subtotal:           subtotal,
				Image:              imageResp,
				CustomDesign:       item.CustomDesign,
			})
			total += subtotal
			continue
		}

		if result.Products == nil {
			continue
		}

		productData, ok := result.Products[*item.ProductID]
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
			CartItemID:         item.ID,
			ProductVariantType: string(item.ProductVariantType),
			ProductID:          item.ProductID,
			ShopID:             item.ShopID,
			ShopName:           shopName,
			ShopSlug:           shopSlug,
			Name:               productData.Product.Name,
			Price:              price,
			Quantity:           quantity,
			Subtotal:           subtotal,
			Image:              image,
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
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	if req.ProductVariantType == "custom" {
		if len(req.CustomDesign) == 0 {
			return apperrors.NewBadRequest("custom_design is required for custom items")
		}
		if req.PhysicalSizeID == "" {
			return apperrors.NewBadRequest("physical_size_id is required for custom items")
		}

		input := usecase.AddCustomItemInput{
			CustomerID:     *authCtx.CustomerID,
			ShopID:         shopID,
			Quantity:       req.Quantity,
			ProductName:    req.ProductName,
			PhysicalSizeID: req.PhysicalSizeID,
			CustomDesign:   req.CustomDesign,
		}
		if err := h.addCustomItem.Execute(r.Context(), input); err != nil {
			return err
		}
	} else {
		// Standard path — existing logic preserved
		if req.ProductID == nil {
			return apperrors.NewBadRequest("product_id is required for standard items")
		}
		productID, err := uuid.Parse(*req.ProductID)
		if err != nil {
			return apperrors.NewBadRequest("invalid product id")
		}
		if req.Quantity <= 0 {
			return apperrors.NewBadRequest("invalid quantity")
		}

		input := usecase.AddItemInput{
			CustomerID: *authCtx.CustomerID,
			ProductID:  productID,
			ShopID:     shopID,
			Quantity:   req.Quantity,
		}
		if err := h.addItem.Execute(r.Context(), input); err != nil {
			return err
		}
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
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
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
		CustomerID: *authCtx.CustomerID,
		ProductID:  productID,
		ShopID:     shopID,
		Quantity:   req.Quantity,
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
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
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
		CustomerID: *authCtx.CustomerID,
		ProductID:  productID,
		ShopID:     shopID,
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

func (h *CartHandler) ChangeItemShop(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	cartItemID, err := apphttp.ParamUUID(r, "cartItemID")
	if err != nil {
		return apperrors.NewBadRequest("invalid cart item id")
	}

	var req changeItemShopRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	newShopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return apperrors.NewBadRequest("invalid target shop id")
	}

	input := usecase.ChangeItemShopInput{
		CustomerID: *authCtx.CustomerID,
		CartItemID: cartItemID,
		NewShopID:  newShopID,
	}

	if err := h.changeItemShop.Execute(r.Context(), input); err != nil {
		return err
	}

	response := map[string]string{
		"message": "item shop updated",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	var req checkoutRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	reqCheckoutCalc := checkoutCalculateRequest{
		Shops: req.Shops,
	}
	input, err := h.parseCheckoutInput(reqCheckoutCalc)
	if err != nil {
		return err
	}

	result, err := h.checkout.Execute(r.Context(), *authCtx, input)
	if err != nil {
		return err
	}

	var shopsResponse []shopResponse
	for _, shop := range result.Shops {
		var itemsResponse []checkoutItemResponse
		for _, item := range shop.Items {
			variantType := "standard"
			if item.IsCustom || item.ProductID == nil {
				variantType = "custom"
			}
			itemsResponse = append(itemsResponse, checkoutItemResponse{
				ProductID:          item.ProductID,
				CartItemID:         item.CartItemID,
				ProductVariantType: variantType,
				IsCustom:           item.IsCustom || item.ProductID == nil,
				ShopID:             item.ShopID,
				Name:               item.Name,
				Price:              item.Price,
				Quantity:           item.Quantity,
				Subtotal:           item.Subtotal,
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
			ShopSlug:     shop.ShopSlug,
			ShopName:     shop.ShopName,
			Subtotal:     shop.Subtotal,
			Total:        &shop.Total,
			Items:        itemsResponse,
			CostCouriers: shippingResponse,
		}

		shopsResponse = append(shopsResponse, shopReponse)
	}

	var paymentMethods []paymentMethodResponse
	for _, pM := range result.PaymentMethods {
		paymentMethod := paymentMethodResponse{
			ID:          pM.ID,
			Name:        pM.Name,
			Type:        pM.Type,
			Description: pM.Description,
			Fee:         pM.Fee,
			Subtotal:    pM.Subtotal,
			Total:       pM.Total,
		}

		paymentMethods = append(paymentMethods, paymentMethod)
	}

	resp := checkoutResponse{
		Address: checkoutAddressResponse{
			ID:            result.Address.ID,
			RecipientName: result.Address.RecipientName,
			Phone:         result.Address.Phone,
			FullAddress:   result.Address.FullAddress,
		},
		Shops:          shopsResponse,
		TotalShipping:  result.TotalShippingFee,
		Subtotal:       result.Subtotal,
		PaymentMethods: paymentMethods,
	}

	if result.TotalAll > 0 {
		resp.TotalAll = &result.TotalAll
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *CartHandler) CheckoutEstimate(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	var req checkoutCalculateRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.AddressID == nil {
		return apperrors.NewBadRequest("address id required")
	}
	if req.PaymentMethodID == nil {
		return apperrors.NewBadRequest("payment method id required")
	}

	for _, shop := range req.Shops {
		if shop.Courier == nil {
			return apperrors.NewBadRequest("courier required")
		}
	}

	input, err := h.parseCheckoutInput(req)
	if err != nil {
		return err
	}

	result, err := h.checkout.Execute(r.Context(), *authCtx, input)
	if err != nil {
		return err
	}

	var shopsResponse []shopCalculateResponse
	for _, shop := range result.Shops {
		var itemsResponse []checkoutItemResponse
		for _, item := range shop.Items {
			variantType := "standard"
			if item.IsCustom || item.ProductID == nil {
				variantType = "custom"
			}
			itemsResponse = append(itemsResponse, checkoutItemResponse{
				ProductID:          item.ProductID,
				CartItemID:         item.CartItemID,
				ProductVariantType: variantType,
				IsCustom:           item.IsCustom || item.ProductID == nil,
				ShopID:             item.ShopID,
				Name:               item.Name,
				Price:              item.Price,
				Quantity:           item.Quantity,
				Subtotal:           item.Subtotal,
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

		shopReponse := shopCalculateResponse{
			ShopID:   shop.ShopID,
			ShopSlug: shop.ShopSlug,
			ShopName: shop.ShopName,
			Subtotal: shop.Subtotal,
			Total:    &shop.Total,
			Items:    itemsResponse,
		}
		if shop.SelectedCourier != nil {
			shopReponse.SelectedCourier = selectedCourierResponse{
				Code:    shop.SelectedCourier.Code,
				Service: shop.SelectedCourier.Service,
				Fee:     shop.SelectedCourier.Fee,
			}
		}

		shopsResponse = append(shopsResponse, shopReponse)
	}

	var paymentMethods []paymentMethodResponse
	for _, pM := range result.PaymentMethods {
		paymentMethod := paymentMethodResponse{
			ID:          pM.ID,
			Name:        pM.Name,
			Type:        pM.Type,
			Description: pM.Description,
			Fee:         pM.Fee,
			Subtotal:    pM.Subtotal,
			Total:       pM.Total,
		}

		paymentMethods = append(paymentMethods, paymentMethod)
	}

	var selectedPayment *paymentMethodResponse
	if result.SelectedPaymentMethod != nil {
		selectedPayment = &paymentMethodResponse{
			ID:          result.SelectedPaymentMethod.ID,
			Name:        result.SelectedPaymentMethod.Name,
			Type:        result.SelectedPaymentMethod.Type,
			Description: result.SelectedPaymentMethod.Description,
			Fee:         result.SelectedPaymentMethod.Fee,
			Subtotal:    result.SelectedPaymentMethod.Subtotal,
			Total:       result.SelectedPaymentMethod.Total,
		}
	}

	resp := checkoutCalculateResponse{
		Address: checkoutAddressResponse{
			ID:            result.Address.ID,
			RecipientName: result.Address.RecipientName,
			Phone:         result.Address.Phone,
			FullAddress:   result.Address.FullAddress,
		},
		Shops:                  shopsResponse,
		TotalShipping:          result.TotalShippingFee,
		Subtotal:               result.Subtotal,
		SelectedPaymentMethods: *selectedPayment,
	}

	if result.TotalAll > 0 {
		resp.TotalAll = &result.TotalAll
	}

	apphttp.WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (h *CartHandler) parseCheckoutInput(
	req checkoutCalculateRequest,
) (usecase.CheckoutInput, error) {
	var paymentMethodID *uuid.UUID
	if req.PaymentMethodID != nil {
		parsed, err := uuid.Parse(*req.PaymentMethodID)
		if err != nil {
			return usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid address id")
		}
		paymentMethodID = &parsed
	}

	var addressID *uuid.UUID
	if req.AddressID != nil {
		parsed, err := uuid.Parse(*req.AddressID)
		if err != nil {
			return usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid address id")
		}
		addressID = &parsed
	}

	var shopInput []usecase.CheckoutShopInput
	for _, shopReq := range req.Shops {
		shopID, err := uuid.Parse(shopReq.ShopID)
		if err != nil {
			return usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid shop id")
		}

		var items []usecase.CheckoutItemInput
		for _, itemReq := range shopReq.Items {
			var productID *uuid.UUID
			var cartItemID *uuid.UUID

			isCustom := (itemReq.IsCustom != nil && *itemReq.IsCustom) ||
				itemReq.ProductVariantType == "custom" ||
				itemReq.ItemType == "custom" ||
				len(itemReq.CustomDesign) > 0 ||
				(itemReq.ProductID == nil && itemReq.CartItemID != nil) ||
				itemReq.ProductID == nil

			if itemReq.ProductID != nil &&
				*itemReq.ProductID != "" &&
				*itemReq.ProductID != "null" &&
				*itemReq.ProductID != "undefined" &&
				*itemReq.ProductID != "custom" {

				if parsed, err := uuid.Parse(*itemReq.ProductID); err == nil {
					productID = &parsed
				} else if !isCustom {
					return usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid product id")
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
				return usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid product id")
			}

			if itemReq.Quantity <= 0 {
				return usecase.CheckoutInput{}, apperrors.NewBadRequest("invalid quantity")
			}

			items = append(items, usecase.CheckoutItemInput{
				ProductID:    productID,
				CartItemID:   cartItemID,
				IsCustom:     isCustom,
				CustomDesign: itemReq.CustomDesign,
				Quantity:     itemReq.Quantity,
			})
		}

		var courier *usecase.SelectedCourierInput
		if shopReq.Courier != nil {
			courier = &usecase.SelectedCourierInput{
				Code:    shopReq.Courier.Code,
				Service: shopReq.Courier.Service,
			}
		}

		shopInput = append(shopInput, usecase.CheckoutShopInput{
			ShopID:  shopID,
			Items:   items,
			Courier: courier,
		})
	}

	return usecase.CheckoutInput{
		PaymentMethodID: paymentMethodID,
		AddressID:       addressID,
		ShopInput:       shopInput,
	}, nil
}

func (h *CartHandler) RemoveCustomItem(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	cartItemID, err := apphttp.ParamUUID(r, "cartItemID")
	if err != nil {
		return apperrors.NewBadRequest("invalid cart item id")
	}

	input := usecase.RemoveCustomItemInput{
		CustomerID: *authCtx.CustomerID,
		CartItemID: cartItemID,
	}

	if err := h.removeCustomItem.Execute(r.Context(), input); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{"message": "item removed"})
	return nil
}
