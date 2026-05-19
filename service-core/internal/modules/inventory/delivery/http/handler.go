package http

import (
	"encoding/json"
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/inventory/usecase"
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

func (h *InventoryHandler) AddInventory(w http.ResponseWriter, r *http.Request) error {
	var req createInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	shopID, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	productID, err := apphttp.ParamUUID(r, "productID")
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	if req.Stock < 0 {
		return apperrors.NewBadRequest("invalid stock")
	}

	if err := h.createInventory.Execute(usecase.CreateInventoryInput{
		ProductID: productID,
		ShopID:    shopID,
		Stock:     req.Stock,
	}); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "inventory successfully added",
	})
	return nil
}
