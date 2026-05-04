package http

import (
	"encoding/json"
	"net/http"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/inventory/usecase"

	"github.com/google/uuid"
)

type InventoryHandler struct {
	createInventory *usecase.CreateInventoryUsecase
}

func NewInventoryHandler(
	createInventory *usecase.CreateInventoryUsecase,
) *InventoryHandler {
	return &InventoryHandler{
		createInventory: createInventory,
	}
}

func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) error {
	var req CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return errors.ErrBadRequest
	}

	shopID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return errors.ErrBadRequest
	}

	if req.Stock < 0 {
		return errors.ErrBadRequest
	}

	if err := h.createInventory.Execute(usecase.CreateInventoryInput{
		ProductID: productID,
		ShopID:    shopID,
		Stock:     req.Stock,
	}); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "inventory successfully created",
	})
	return nil
}
