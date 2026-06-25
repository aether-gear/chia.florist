package shipping

import (
	"context"

	locationDomain "service-core/internal/modules/location/domain"
)

type CalculateRatesInput struct {
	OriginID      int
	DestinationID int
	Weight        int
	Couriers      []string
	PriceFilter   *string
}

type RateOption struct {
	Name        string
	Code        string
	Service     string
	Description string
	Cost        int64
	Etd         string
}

type Provider interface {
	CalculateRates(
		ctx context.Context,
		input CalculateRatesInput,
	) ([]RateOption, error)

	// TrackShipment(
	// 	ctx context.Context,
	// 	trackingNumber string,
	// ) ([]shipmentDomain.ShipmentEvent, error)

	ListProvinces(
		ctx context.Context,
	) ([]locationDomain.Province, error)

	ListCities(
		ctx context.Context,
		provinceID string,
	) ([]locationDomain.City, error)

	ListDistricts(
		ctx context.Context,
		cityID string,
	) ([]locationDomain.District, error)

	ListVillages(
		ctx context.Context,
		districtID string,
	) ([]locationDomain.Village, error)
}
