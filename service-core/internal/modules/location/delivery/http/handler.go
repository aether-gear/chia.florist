package http

import (
	"net/http"

	"service-core/internal/modules/location/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
)

type LocationHandler struct {
	listLocation *usecase.ListLocationUsecase
}

func NewLocationHandler(
	listLocation *usecase.ListLocationUsecase,
) *LocationHandler {
	return &LocationHandler{
		listLocation: listLocation,
	}
}

func (h *LocationHandler) Province(w http.ResponseWriter, r *http.Request) error {
	result, err := h.listLocation.Province()
	if err != nil {
		return err
	}

	res := make([]ProvinceResponse, 0, len(result))
	for _, r := range result {
		province := ProvinceResponse{
			ID:   r.ID,
			Name: r.Name,
		}

		res = append(res, province)
	}

	apphttp.WriteJSON(w, http.StatusOK, res)
	return nil
}

func (h *LocationHandler) City(w http.ResponseWriter, r *http.Request) error {
	provinceID := r.URL.Query().Get("province_id")
	if provinceID == "" {
		return errors.ErrBadRequest
	}

	result, err := h.listLocation.City(provinceID)
	if err != nil {
		return err
	}

	res := make([]CityResponse, 0, len(result))
	for _, r := range result {
		city := CityResponse{
			ID:         r.ID,
			ProvinceID: r.ProvinceID,
			Name:       r.Name,
		}

		res = append(res, city)
	}

	apphttp.WriteJSON(w, http.StatusOK, res)
	return nil
}

func (h *LocationHandler) District(w http.ResponseWriter, r *http.Request) error {
	cityID := r.URL.Query().Get("city_id")
	if cityID == "" {
		return errors.ErrBadRequest
	}

	result, err := h.listLocation.District(cityID)
	if err != nil {
		return err
	}

	res := make([]DistrictResponse, 0, len(result))
	for _, r := range result {
		district := DistrictResponse{
			ID:     r.ID,
			CityID: r.CityID,
			Name:   r.Name,
		}

		res = append(res, district)
	}

	apphttp.WriteJSON(w, http.StatusOK, res)
	return nil
}

func (h *LocationHandler) Village(w http.ResponseWriter, r *http.Request) error {
	districtID := r.URL.Query().Get("district_id")
	if districtID == "" {
		return errors.ErrBadRequest
	}

	result, err := h.listLocation.Village(districtID)
	if err != nil {
		return err
	}

	res := make([]VillageResponse, 0, len(result))
	for _, r := range result {
		village := VillageResponse{
			ID:         r.ID,
			DistrictID: r.DistrictID,
			Name:       r.Name,
		}

		res = append(res, village)
	}

	apphttp.WriteJSON(w, http.StatusOK, res)
	return nil
}
