package domain

type AddressDetail struct {
	ProvinceID  string
	CityID      string
	DistrictID  string
	VillageID   string
	FullAddress string
	PostalCode  string

	Latitude  *float64
	Longitude *float64
}
