package http

type ProvinceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CityResponse struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

type DistrictResponse struct {
	ID     string `json:"id"`
	CityID string `json:"city_id"`
	Name   string `json:"name"`
}

type VillageResponse struct {
	ID         string `json:"id"`
	DistrictID string `json:"district_id"`
	Name       string `json:"name"`
}
