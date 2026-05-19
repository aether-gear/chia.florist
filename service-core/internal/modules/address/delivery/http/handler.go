package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/address/usecase"

	"github.com/google/uuid"
)

type AddressHandler struct {
	getShopAddress    *usecase.GetShopAddressUsecase
	listUserAddresses *usecase.ListUserAddressUsecase
	listShopAddresses *usecase.ListShopAddressesUsecase
	createUserAddress *usecase.CreateAddressUsecase
	createShopAddress *usecase.CreateShopAddressUsecase
}

func NewAddressHandler(
	getShopAddress *usecase.GetShopAddressUsecase,
	listUserAddresses *usecase.ListUserAddressUsecase,
	listShopAddresses *usecase.ListShopAddressesUsecase,
	createUserAddress *usecase.CreateAddressUsecase,
	createShopAddress *usecase.CreateShopAddressUsecase,
) *AddressHandler {
	return &AddressHandler{
		getShopAddress:    getShopAddress,
		listUserAddresses: listUserAddresses,
		listShopAddresses: listShopAddresses,
		createUserAddress: createUserAddress,
		createShopAddress: createShopAddress,
	}
}

func (h *AddressHandler) ListUserAddresses(w http.ResponseWriter, r *http.Request) error {
	userID, err := apphttp.ParamUUID(r, "userID")
	if err != nil {
		return apperrors.NewBadRequest("invallid user id")
	}

	result, err := h.listUserAddresses.ListByUserID(userID)
	if err != nil {
		return err
	}

	response := make([]userAddressResponse, 0, len(result))
	for _, r := range result {
		address := userAddressResponse{
			UserID:       r.UserID,
			ReceiverName: r.ReceiverName,
			Phone:        r.Phone,
			IsDefault:    r.IsDefault,
			ProvinceID:   r.Detail.ProvinceID,
			CityID:       r.Detail.CityID,
			DistrictID:   r.Detail.DistrictID,
			VillageID:    r.Detail.VillageID,
			FullAddress:  r.Detail.FullAddress,
			PostalCode:   r.Detail.PostalCode,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}

		response = append(response, address)
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) CreateUserAddress(w http.ResponseWriter, r *http.Request) error {
	var req createUserAddressRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	parsedUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		return apperrors.NewBadRequest("invalid user id")
	}

	var parsedIsActive bool
	parsedIsActive, err = strconv.ParseBool(*req.IsDefault)
	if err != nil {
		return apperrors.NewBadRequest("invalid active status")
	}

	input := usecase.CreateAddressInput{
		UserID:       parsedUserID,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		IsDefault:    &parsedIsActive,
		ProvinceID:   req.ProvinceID,
		CityID:       req.CityID,
		DistrictID:   req.DistrictID,
		VillageID:    req.VillageID,
		FullAddress:  req.FullAddress,
		PostalCode:   req.PostalCode,
	}

	err = h.createUserAddress.Execute(input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "address successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) GetShopAddress(w http.ResponseWriter, r *http.Request) error {
	addressID, err := apphttp.ParamUUID(r, "addressID")
	if err != nil {
		return apperrors.NewBadRequest("invalid address id")
	}

	result, err := h.getShopAddress.GetByID(addressID)
	if err != nil {
		return err
	}

	response := shopAddressResponse{
		ShopID:      result.ShopID,
		Label:       result.Label,
		Phone:       result.Phone,
		IsActive:    result.IsActive,
		ProvinceID:  result.Detail.ProvinceID,
		CityID:      result.Detail.CityID,
		DistrictID:  result.Detail.DistrictID,
		VillageID:   result.Detail.VillageID,
		FullAddress: result.Detail.FullAddress,
		PostalCode:  result.Detail.PostalCode,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) ListShopAddresses(w http.ResponseWriter, r *http.Request) error {
	shopID, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid shop id")
	}

	result, err := h.listShopAddresses.FindByShopID(shopID)
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	parsedShopID, err := uuid.Parse(req.ShopID)
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

	err = h.createShopAddress.Execute(input)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "address successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
