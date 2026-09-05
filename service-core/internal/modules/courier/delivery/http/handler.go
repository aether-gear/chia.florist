package http

import (
	"fmt"
	"net/http"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	"service-core/internal/modules/courier/domain"
	"service-core/internal/modules/courier/usecase"

	"github.com/google/uuid"
)

type CourierHandler struct {
	listCouriers         *usecase.ListCouriersUsecase
	configureShopCourier *usecase.ConfigureShopCourierUsecase
	verifyShopCourier    *usecase.VerifyShopCourierUsecase
}

func NewCourierHandler(
	listCouriers *usecase.ListCouriersUsecase,
	configureShopCourier *usecase.ConfigureShopCourierUsecase,
	verifyShopCourier *usecase.VerifyShopCourierUsecase,
) *CourierHandler {
	return &CourierHandler{
		listCouriers:         listCouriers,
		configureShopCourier: configureShopCourier,
		verifyShopCourier:    verifyShopCourier,
	}
}

func (h *CourierHandler) ListAllCouriers(w http.ResponseWriter, r *http.Request) error {
	code, err := h.listCouriers.Execute(r.Context())
	if err != nil {
		return err
	}

	response := map[string][]string{
		"couriers": code,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CourierHandler) ConfigureCourierShop(w http.ResponseWriter, r *http.Request) error {
	var req configureCourierShopRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	shopID, err := apphttp.ParamUUID(r, "shopID")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	inputs := make([]usecase.ConfigureShopCourierInput, 0, len(req.Couriers))
	for _, courier := range req.Couriers {
		inputs = append(inputs, usecase.ConfigureShopCourierInput{
			Code:   courier.Code,
			Active: courier.IsEnabled,
		})
	}

	err = h.configureShopCourier.
		Execute(r.Context(), shopID, inputs)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "courier shops successfully configured",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *CourierHandler) UpdateShopCourier(w http.ResponseWriter, r *http.Request) error {
	var req updateShopCourierRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	shopID, err := apphttp.ParamUUID(r, "shopID")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	code := apphttp.Param(r, "code")
	if code == "" {
		return apperrors.NewBadRequest("courier code is required")
	}

	var isAdmin bool
	var staffID *uuid.UUID
	if actor, ok := authzSvc.GetActor(r.Context()); ok {
		isAdmin = actor.IsSuperAdmin()
		staffID = actor.StaffID
	}

	updated, err := h.configureShopCourier.UpdateSingle(r.Context(), usecase.UpdateSingleShopCourierInput{
		ShopID:          shopID,
		Code:            code,
		Name:            req.Name,
		LocationAddress: req.LocationAddress,
		Active:          req.Active,
		IsAdmin:         isAdmin,
		AdminStaffID:    staffID,
	})
	if err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "courier updated successfully",
		"courier": toShopCourierDetailResponse(updated),
	})
	return nil
}

func (h *CourierHandler) VerifyShopCourier(w http.ResponseWriter, r *http.Request) error {
	var req verifyShopCourierRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	shopID, err := apphttp.ParamUUID(r, "shopID")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	code := apphttp.Param(r, "code")
	if code == "" {
		return apperrors.NewBadRequest("courier code is required")
	}

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok || actor.StaffID == nil {
		return apperrors.NewUnauthorized("admin staff authentication required")
	}

	updated, err := h.verifyShopCourier.Execute(r.Context(), usecase.VerifyShopCourierInput{
		ShopID:          shopID,
		Code:            code,
		Action:          req.Action,
		RejectionReason: req.RejectionReason,
		AdminStaffID:    *actor.StaffID,
	})
	if err != nil {
		return err
	}

	apphttp.WriteJSON(w, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("courier %s successfully", req.Action),
		"courier": toShopCourierDetailResponse(updated),
	})
	return nil
}

func toShopCourierDetailResponse(c *domain.ShopCourier) shopCourierDetailResponse {
	var verifiedAtStr, verifiedByStr, rejectionReasonStr *string
	if c.VerifiedAt != nil {
		str := c.VerifiedAt.Format(time.RFC3339)
		verifiedAtStr = &str
	}
	if c.VerifiedBy != nil {
		str := c.VerifiedBy.String()
		verifiedByStr = &str
	}
	if c.RejectionReason != nil {
		rejectionReasonStr = c.RejectionReason
	}

	return shopCourierDetailResponse{
		ShopID:             c.ShopID.String(),
		Code:               c.Code,
		BranchName:         c.BranchName,
		Name:               c.Name,
		LocationAddress:    c.LocationAddress,
		Active:             c.Active,
		VerificationStatus: string(c.VerificationStatus),
		VerifiedAt:         verifiedAtStr,
		VerifiedBy:         verifiedByStr,
		RejectionReason:    rejectionReasonStr,
	}
}
