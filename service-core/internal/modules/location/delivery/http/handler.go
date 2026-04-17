package http

import (
	"encoding/json"
	"net/http"
	"service-core/internal/modules/location/usecase"
)

type LocationHandler struct {
	listUC *usecase.ListLocationUsecase
}

func NewLocationHandler(listUC *usecase.ListLocationUsecase) *LocationHandler {
	return &LocationHandler{
		listUC: listUC,
	}
}

func (h *LocationHandler) Province(w http.ResponseWriter, r *http.Request) {
	result, err := h.listUC.Province()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := make([]ProvinceResponse, 0, len(result))
	for _, r := range result {
		province := ProvinceResponse{
			ID:   r.ID,
			Name: r.Name,
		}

		res = append(res, province)
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *LocationHandler) City(w http.ResponseWriter, r *http.Request) {
	provinceID := r.URL.Query().Get("province_id")
	if provinceID == "" {
		http.Error(w, "province_id is required", http.StatusBadRequest)
		return
	}

	result, err := h.listUC.City(provinceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	writeJSON(w, http.StatusOK, res)
}

func (h *LocationHandler) District(w http.ResponseWriter, r *http.Request) {
	cityID := r.URL.Query().Get("city_id")
	if cityID == "" {
		http.Error(w, "city_id is required", http.StatusBadRequest)
		return
	}

	result, err := h.listUC.District(cityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	writeJSON(w, http.StatusOK, res)
}

func (h *LocationHandler) Village(w http.ResponseWriter, r *http.Request) {
	districtID := r.URL.Query().Get("district_id")
	if districtID == "" {
		http.Error(w, "district_id is required", http.StatusBadRequest)
		return
	}

	result, err := h.listUC.Village(districtID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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

	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
