package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"service-core/internal/modules/address/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

	"github.com/google/uuid"
)

type AddressHandler struct {
	getAddress    *usecase.GetAddressUsecase
	createAddress *usecase.CreateAddressUsecase
}

func NewAddressHandler(
	getAddress *usecase.GetAddressUsecase,
	createAddress *usecase.CreateAddressUsecase,
) *AddressHandler {
	return &AddressHandler{
		getAddress:    getAddress,
		createAddress: createAddress,
	}
}

func (h *AddressHandler) GetAddresses(w http.ResponseWriter, r *http.Request) error {
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

	result, err := h.getAddress.GetByUserID(parsedID)
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
			ProvinceID:   r.ProvinceID,
			CityID:       r.CityID,
			DistrictID:   r.DistrictID,
			VillageID:    r.VillageID,
			FullAddress:  r.FullAddress,
			PostalCode:   r.PostalCode,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		}

		response = append(response, address)
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) error {
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

	err = h.createAddress.Execute(inputCreateAddress)
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "address successfully created",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
