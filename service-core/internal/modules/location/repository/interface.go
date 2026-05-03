package repository

import "service-core/internal/modules/location/domain"

type LocationRepository interface {
	ListProvinces() ([]domain.Province, error)
	ListCitiesByProvince(provinceID string) ([]domain.City, error)
	ListDistrictsByCity(cityID string) ([]domain.District, error)
	ListVillagesByDistrict(districtID string) ([]domain.Village, error)
}
