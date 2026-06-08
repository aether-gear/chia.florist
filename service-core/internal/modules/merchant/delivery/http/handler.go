package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authenDomain "service-core/internal/modules/authentication/domain"
	authzSvc "service-core/internal/modules/authorization/infra/service"
	"service-core/internal/modules/merchant/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type merchantHandler struct {
	addMerchantAccount *usecase.AddMerchantAccountUsecase
	createMerchant     *usecase.CreateMerchantUsecase
}

func NewMerchantHandler(
	addMerchantAccount *usecase.AddMerchantAccountUsecase,
	createMerchant *usecase.CreateMerchantUsecase,
) *merchantHandler {
	return &merchantHandler{
		addMerchantAccount: addMerchantAccount,
		createMerchant:     createMerchant,
	}
}

func (h *merchantHandler) AddMerchantAccount(w http.ResponseWriter, r *http.Request) error {
	merchantIDStr := chi.URLParam(r, "merchantID")
	merchantID, err := uuid.Parse(merchantIDStr)
	if err != nil {
		return apperrors.NewBadRequest("invalid merchant id")
	}

	authCtx, ok := authenDomain.GetAuthContext(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}
	_ = authCtx

	actor, ok := authzSvc.GetActor(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}

	var req addMerchantAccountRequest
	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("email is required")
	}
	if req.Name == "" {
		return apperrors.NewBadRequest("name is required")
	}
	if req.Username == "" {
		return apperrors.NewBadRequest("username is required")
	}

	input := usecase.AddMerchantAccountParams{
		ActorAccountID:  actor.AccountID,
		ActorMerchantID: *actor.MerchantID,
		MerchantID:      merchantID,
		Email:           req.Email,
		Name:            req.Name,
		Password:        req.Password,
		Username:        req.Username,
		Phone:           req.Phone,
	}

	err = h.addMerchantAccount.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "verify success",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *merchantHandler) CreateMerchant(w http.ResponseWriter, r *http.Request) error {
	var req createMerchantRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.Name == "" {
		return apperrors.NewBadRequest("invalid name")
	}
	if req.Description != nil && *req.Description == "" {
		return apperrors.NewBadRequest("invalid description")
	}

	if req.LogoUrl != nil && *req.LogoUrl == "" {
		return apperrors.NewBadRequest("invalid logo url")
	}

	if req.BannerUrl != nil && *req.BannerUrl == "" {
		return apperrors.NewBadRequest("invalid banner url")
	}

	input := usecase.CreateMerchantInput{
		Name:        req.Name,
		Description: req.Description,
		LogoUrl:     req.LogoUrl,
		BannerUrl:   req.BannerUrl,
	}

	err := h.createMerchant.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "merchant successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
