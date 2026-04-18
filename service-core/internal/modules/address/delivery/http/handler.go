package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"service-core/internal/modules/address/usecase"
	"strconv"
	"strings"

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

func (h *AddressHandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[3] == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id := parts[3]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	result, err := h.getAddress.GetByUserID(parsedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	json.NewEncoder(w).Encode(response)
}

func (h *AddressHandler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	var req CreateAddressRequest

	body, _ := io.ReadAll(r.Body)

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProvinceID == "" || req.CityID == "" || req.DistrictID == "" || req.VillageID == "" {
		http.Error(w, "missing some locations", http.StatusBadRequest)
		return
	}

	if req.FullAddress == "" {
		http.Error(w, "missing full address", http.StatusBadRequest)
		return
	}

	if req.PostalCode == "" {
		http.Error(w, "missing postal code", http.StatusBadRequest)
		return
	}

	parsedID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var parsedBool bool
	parsedBool, err = strconv.ParseBool(*req.IsDefault)
	if err != nil {
		http.Error(w, "invalid is default value", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "address successfully created",
	})
}
