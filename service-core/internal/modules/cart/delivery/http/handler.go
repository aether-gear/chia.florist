package http

import (
	"encoding/json"
	"net/http"
	"service-core/internal/modules/cart/usecase"

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

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	result, err := h.getCart.Execute(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result == nil || result.Cart == nil {
		http.Error(w, "cart not found", http.StatusNotFound)
		return
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

	json.NewEncoder(w).Encode(response)
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req addItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		http.Error(w, "invalid product_id", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	if err := h.addItem.Execute(userID, productID, req.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "item added",
	})
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var req updateItemRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		http.Error(w, "invalid product_id", http.StatusBadRequest)
		return
	}

	if req.Quantity < 0 {
		http.Error(w, "quantity cannot be negative", http.StatusBadRequest)
		return
	}

	if err := h.updateItem.Execute(userID, productID, req.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "item updated",
	})
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	productIDStr := r.URL.Query().Get("product_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "invalid product_id", http.StatusBadRequest)
		return
	}

	if err := h.removeItem.Execute(userID, productID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "item removed",
	})
}
