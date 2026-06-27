package http

import (
	"net/http"
	"strconv"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/address/usecase"
	authendomain "service-core/internal/modules/authentication/domain"

	"github.com/google/uuid"
)

type AddressHandler struct {
	listCustomerAddresses *usecase.ListCustomerAddressesUsecase
	saveCustomerAddress   *usecase.SaveCustomerAddressUsecase
	deleteCustomerAddress *usecase.DeleteCustomerAddressUsecase
	listShopAddresses     *usecase.ListShopAddressesUsecase
	createShopAddress     *usecase.CreateShopAddressUsecase
}

func NewAddressHandler(
	listCustomerAddresses *usecase.ListCustomerAddressesUsecase,
	saveCustomerAddress *usecase.SaveCustomerAddressUsecase,
	deleteCustomerAddress *usecase.DeleteCustomerAddressUsecase,
	listShopAddresses *usecase.ListShopAddressesUsecase,
	createShopAddress *usecase.CreateShopAddressUsecase,
) *AddressHandler {
	return &AddressHandler{
		listCustomerAddresses: listCustomerAddresses,
		saveCustomerAddress:   saveCustomerAddress,
		deleteCustomerAddress: deleteCustomerAddress,
		listShopAddresses:     listShopAddresses,
		createShopAddress:     createShopAddress,
	}
}

func (h *AddressHandler) ListUserAddresses(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authendomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	addresses, err := h.listCustomerAddresses.ListByCustomerID(r.Context(), *authCtx.CustomerID)
	if err != nil {
		return err
	}

	result := make([]customerAddressResponse, 0, len(addresses))
	for _, r := range addresses {
		address := customerAddressResponse{
			AddressID:    r.ID,
			CustomerID:   r.CustomerID,
			ReceiverName: r.ReceiverName,
			Phone:        r.Phone,
			IsDefault:    r.IsDefault,
			ProvinceID:   r.Detail.ProvinceID,
			CityID:       r.Detail.CityID,
			DistrictID:   r.Detail.DistrictID,
			VillageID:    r.Detail.VillageID,
			FullAddress:  r.Detail.FullAddress,
			PostalCode:   r.Detail.PostalCode,
			// CreatedAt:    r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}

		result = append(result, address)
	}

	response := map[string][]customerAddressResponse{
		"addresses": result,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) SaveUserAddress(w http.ResponseWriter, r *http.Request) error {
	var req saveCustomerAddressRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.ProvinceID == "" {
		return apperrors.NewBadRequest("invalid province id")
	}
	if req.DistrictID == "" {
		return apperrors.NewBadRequest("invalid district id")
	}
	if req.CityID == "" {
		return apperrors.NewBadRequest("invalid city id")
	}
	if req.VillageID == "" {
		return apperrors.NewBadRequest("invalid village id")
	}
	if req.FullAddress == "" {
		return apperrors.NewBadRequest("invalid full address")
	}
	if req.PostalCode == "" {
		return apperrors.NewBadRequest("invalid postal code")
	}

	authCtx, ok := authendomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	var addressID *uuid.UUID
	if req.AddressID != nil {
		parsed, err := uuid.Parse(*req.AddressID)
		if err != nil {
			return apperrors.NewBadRequest("invalid address id")
		}

		addressID = &parsed
	}

	var parsedIsDefault = false
	if req.IsDefault != nil && *req.IsDefault != "" {
		parsed, err := strconv.ParseBool(*req.IsDefault)
		if err != nil {
			return apperrors.NewBadRequest("invalid default status")
		}
		parsedIsDefault = parsed
	}

	input := usecase.SaveCustomerAddressInput{
		ID:           addressID,
		CustomerID:   *authCtx.CustomerID,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		IsDefault:    &parsedIsDefault,
		ProvinceID:   req.ProvinceID,
		CityID:       req.CityID,
		DistrictID:   req.DistrictID,
		VillageID:    req.VillageID,
		FullAddress:  req.FullAddress,
		PostalCode:   req.PostalCode,
	}

	err := h.saveCustomerAddress.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "address saved successfully",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) DeleteUserAddress(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authendomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}
	if authCtx.CustomerID == nil {
		return apperrors.NewForbidden("customer account required")
	}

	addressID, err := apphttp.ParamUUID(r, "addressID")
	if err != nil {
		return apperrors.NewBadRequest("invalid address id")
	}

	err = h.deleteCustomerAddress.Execute(r.Context(), addressID)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "address deleted successfully",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) ListShopAddresses(w http.ResponseWriter, r *http.Request) error {
	shopID, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	result, err := h.listShopAddresses.FindByShopID(r.Context(), shopID)
	if err != nil {
		return err
	}

	addresses := make([]shopAddressResponse, 0, len(result))
	for _, r := range result {
		address := shopAddressResponse{
			ShopID:      r.ShopID,
			Label:       r.Label,
			Phone:       r.Phone,
			IsActive:    r.IsActive,
			ProvinceID:  r.Detail.ProvinceID,
			CityID:      r.Detail.CityID,
			DistrictID:  r.Detail.DistrictID,
			VillageID:   r.Detail.VillageID,
			FullAddress: r.Detail.FullAddress,
			PostalCode:  r.Detail.PostalCode,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}

		addresses = append(addresses, address)
	}

	response := map[string][]shopAddressResponse{
		"addresses": addresses,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) CreateShopAddress(w http.ResponseWriter, r *http.Request) error {
	var req createShopAddressRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid request body")
	}

	if req.ProvinceID == "" {
		return apperrors.NewBadRequest("invalid province id")
	}
	if req.DistrictID == "" {
		return apperrors.NewBadRequest("invalid district id")
	}
	if req.CityID == "" {
		return apperrors.NewBadRequest("invalid city id")
	}
	if req.VillageID == "" {
		return apperrors.NewBadRequest("invalid village id")
	}

	if req.FullAddress == "" {
		return apperrors.NewBadRequest("invalid full address")
	}
	if req.PostalCode == "" {
		return apperrors.NewBadRequest("invalid postal code")
	}

	parsedShopID, err := apphttp.ParamUUID(r, "shopID")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	var parsedIsActive bool
	parsedIsActive, err = strconv.ParseBool(req.IsActive)
	if err != nil {
		return apperrors.NewBadRequest("invalid active status")
	}

	input := usecase.CreateShopAddressInput{
		ShopID:      parsedShopID,
		Label:       req.Label,
		Phone:       req.Phone,
		IsActive:    &parsedIsActive,
		ProvinceID:  req.ProvinceID,
		CityID:      req.CityID,
		DistrictID:  req.DistrictID,
		VillageID:   req.VillageID,
		FullAddress: req.FullAddress,
		PostalCode:  req.PostalCode,
	}

	err = h.createShopAddress.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "address successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
