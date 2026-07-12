package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/inventory/usecase"
)

type InventoryHandler struct {
	createInventory *usecase.CreateInventoryUsecase
	updateInventory *usecase.UpdateInventoryUsecase
	deleteInventory *usecase.DeleteInventoryUsecase
}

func NewInventoryHandler(
	createInventory *usecase.CreateInventoryUsecase,
	updateInventory *usecase.UpdateInventoryUsecase,
	deleteInventory *usecase.DeleteInventoryUsecase,
) *InventoryHandler {
	return &InventoryHandler{
		createInventory: createInventory,
		updateInventory: updateInventory,
		deleteInventory: deleteInventory,
	}
}

func (h *InventoryHandler) AddInventory(w http.ResponseWriter, r *http.Request) error {
	var req createInventoryRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	shopID, err := apphttp.ParamUUID(r, "shopID")
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

	if err := h.createInventory.Execute(
		r.Context(),
		usecase.CreateInventoryInput{
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

func (h *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) error {
	var req updateInventoryRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	shopID, err := apphttp.ParamUUID(r, "shopID")
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

	if err := h.updateInventory.Execute(
		r.Context(),
		usecase.UpdateInventoryInput{
			ProductID: productID,
			ShopID:    shopID,
			Stock:     req.Stock,
		}); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "inventory successfully updated",
	})
	return nil
}

func (h *InventoryHandler) RemoveInventory(w http.ResponseWriter, r *http.Request) error {
	shopID, err := apphttp.ParamUUID(r, "shopID")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	productID, err := apphttp.ParamUUID(r, "productID")
	if err != nil {
		return apperrors.NewBadRequest("invalid product id")
	}

	if err := h.deleteInventory.Execute(
		r.Context(),
		usecase.DeleteInventoryInput{
			ProductID: productID,
			ShopID:    shopID,
		}); err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "inventory successfully removed",
	})
	return nil
}
