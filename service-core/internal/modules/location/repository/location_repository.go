package repository

import "github.com/google/uuid"

type LocationRepository interface {
	GetProvinces() error
	GetCitiesByProvince(provinceID uuid.UUID) error
	GetDistrictsByCity(cityID uuid.UUID) error
}
