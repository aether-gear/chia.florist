package http

type provinceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type cityResponse struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

type districtResponse struct {
	ID     string `json:"id"`
	CityID string `json:"city_id"`
	Name   string `json:"name"`
}

type villageResponse struct {
	ID         string `json:"id"`
	DistrictID string `json:"district_id"`
	Name       string `json:"name"`
}
