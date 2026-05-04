package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"service-core/internal/common/errors"
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
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[3] == "" {
		return errors.ErrBadRequest
	}

	id := parts[3]
	if id == "" {
		return errors.ErrBadRequest
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	result, err := h.listUserAddresses.ListByUserID(parsedID)
	if err != nil {
		return err
	}

	response := make([]AddressResponse, 0, len(result))
	for _, r := range result {
		address := AddressResponse{
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
	var req CreateAddressRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.ProvinceID == "" || req.CityID == "" || req.DistrictID == "" || req.VillageID == "" {
		return errors.ErrBadRequest
	}

	if req.FullAddress == "" || req.PostalCode == "" {
		return errors.ErrBadRequest
	}

	parsedID, err := uuid.Parse(req.UserID)
	if err != nil {
		return errors.ErrBadRequest
	}

	var parsedBool bool
	parsedBool, err = strconv.ParseBool(*req.IsDefault)
	if err != nil {
		return errors.ErrBadRequest
	}

	inputCreateAddress := usecase.CreateAddressInput{
		UserID:       parsedID,
		ReceiverName: req.ReceiverName,
		Phone:        req.Phone,
		IsDefault:    &parsedBool,
		ProvinceID:   req.ProvinceID,
		CityID:       req.CityID,
		DistrictID:   req.DistrictID,
		VillageID:    req.VillageID,
		FullAddress:  req.FullAddress,
		PostalCode:   req.PostalCode,
	}

	err = h.createUserAddress.Execute(inputCreateAddress)
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
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[3] == "" {
		return errors.ErrBadRequest
	}

	addressID := parts[3]
	if addressID == "" {
		return errors.ErrBadRequest
	}

	parsedAddressID, err := uuid.Parse(addressID)
	if err != nil {
		return errors.ErrBadRequest
	}

	result, err := h.getShopAddress.GetByID(parsedAddressID)
	if err != nil {
		return err
	}

	response := ShopAddressResponse{
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
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[3] == "" {
		return errors.ErrBadRequest
	}

	shopID := parts[3]
	if shopID == "" {
		return errors.ErrBadRequest
	}

	parsedShopID, err := uuid.Parse(shopID)
	if err != nil {
		return errors.ErrBadRequest
	}

	result, err := h.listShopAddresses.FindByShopID(parsedShopID)
	if err != nil {
		return err
	}

	addresses := make([]ShopAddressResponse, 0, len(result))
	for _, r := range result {
		address := ShopAddressResponse{
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

	response := ShopAddressesResponse{
		Addresses: addresses,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) CreateShopAddress(w http.ResponseWriter, r *http.Request) error {
	var req CreateShopAddressRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.ProvinceID == "" || req.CityID == "" || req.DistrictID == "" || req.VillageID == "" {
		return errors.ErrBadRequest
	}

	if req.FullAddress == "" || req.PostalCode == "" {
		return errors.ErrBadRequest
	}

	parsedID, err := uuid.Parse(req.ShopID)
	if err != nil {
		return errors.ErrBadRequest
	}

	var parsedBool bool
	parsedBool, err = strconv.ParseBool(req.IsActive)
	if err != nil {
		return errors.ErrBadRequest
	}

	inputCreateAddress := usecase.CreateShopAddressInput{
		ShopID:      parsedID,
		Label:       req.Label,
		Phone:       req.Phone,
		IsActive:    &parsedBool,
		ProvinceID:  req.ProvinceID,
		CityID:      req.CityID,
		DistrictID:  req.DistrictID,
		VillageID:   req.VillageID,
		FullAddress: req.FullAddress,
		PostalCode:  req.PostalCode,
	}

	err = h.createShopAddress.Execute(inputCreateAddress)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "address successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
