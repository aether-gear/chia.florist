package http

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/location/usecase"
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
	result, err := h.listLocation.Province(r.Context())
	if err != nil {
		return err
	}

	provinces := make([]provinceResponse, 0, len(result))
	for _, r := range result {
		p := provinceResponse{
			ID:   r.ID,
			Name: r.Name,
		}

		provinces = append(provinces, p)
	}

	response := map[string][]provinceResponse{
		"provinces": provinces,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *LocationHandler) City(w http.ResponseWriter, r *http.Request) error {
	provinceID := apphttp.Param(r, "id")

	result, err := h.listLocation.City(r.Context(), provinceID)
	if err != nil {
		return err
	}

	cities := make([]cityResponse, 0, len(result))
	for _, r := range result {
		c := cityResponse{
			ID:         r.ID,
			ProvinceID: r.ProvinceID,
			Name:       r.Name,
		}

		cities = append(cities, c)
	}

	response := map[string][]cityResponse{
		"cities": cities,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *LocationHandler) District(w http.ResponseWriter, r *http.Request) error {
	cityID := apphttp.Param(r, "id")

	result, err := h.listLocation.District(r.Context(), cityID)
	if err != nil {
		return err
	}

	districts := make([]districtResponse, 0, len(result))
	for _, r := range result {
		d := districtResponse{
			ID:     r.ID,
			CityID: r.CityID,
			Name:   r.Name,
		}

		districts = append(districts, d)
	}

	response := map[string][]districtResponse{
		"districts": districts,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *LocationHandler) Village(w http.ResponseWriter, r *http.Request) error {
	districtID := apphttp.Param(r, "id")

	result, err := h.listLocation.Village(r.Context(), districtID)
	if err != nil {
		return err
	}

	villages := make([]villageResponse, 0, len(result))
	for _, r := range result {
		v := villageResponse{
			ID:         r.ID,
			DistrictID: r.DistrictID,
			Name:       r.Name,
		}

		villages = append(villages, v)
	}

	response := map[string][]villageResponse{
		"villages": villages,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
