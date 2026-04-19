package http

import (
	"encoding/json"
	"net/http"

	"service-core/internal/modules/cart/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

	"github.com/google/uuid"
)

type CartHandler struct {
	addItem    *usecase.AddItemUsecase
	getCart    *usecase.GetCartUsecase
	updateItem *usecase.UpdateItemUsecase
	removeItem *usecase.RemoveItemUsecase
}

func NewCartHandler(
	aI *usecase.AddItemUsecase,
	gC *usecase.GetCartUsecase,
	uI *usecase.UpdateItemUsecase,
	rI *usecase.RemoveItemUsecase,
) *CartHandler {
	return &CartHandler{
		addItem:    aI,
		getCart:    gC,
		updateItem: uI,
		removeItem: rI,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) error {
	userID := r.URL.Query().Get("user_id")
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.ErrBadRequest
	}

	result, err := h.getCart.Execute(parsedUserID)
	if err != nil {
		return err
	}
	if result == nil || result.Cart == nil {
		return errors.ErrNotFound
	}

	var total int64
	items := make([]CartItemView, 0, len(result.Cart.Items))

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

		items = append(items, CartItemView{
			ProductID: item.ID,
			Name:      productData.Product.Name,
			Price:     price,
			Quantity:  quantity,
			Subtotal:  subtotal,
		})

		total += subtotal
	}

	response := CartResponse{
		CartID: result.Cart.ID,
		Items:  items,
		Total:  total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) error {
	var req addItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.ErrBadRequest
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return errors.ErrBadRequest
	}

	if req.Quantity <= 0 {
		return errors.ErrBadRequest
	}

	if err := h.addItem.Execute(userID, productID, req.Quantity); err != nil {
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return errors.ErrBadRequest
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.ErrBadRequest
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return errors.ErrBadRequest
	}

	if req.Quantity < 0 {
		return errors.ErrBadRequest
	}

	if err := h.updateItem.Execute(userID, productID, req.Quantity); err != nil {
		return err
	}

	response := map[string]string{
		"message": "item updated",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) error {
	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return errors.ErrBadRequest
	}

	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := h.removeItem.Execute(parsedUserID, parsedProductID); err != nil {
		return err
	}

	response := map[string]string{
		"message": "item removed",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
