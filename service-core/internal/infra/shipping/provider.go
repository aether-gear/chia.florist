package shipping

import (
	"context"
	"time"

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

type ShippingProvider interface {
	CalculateRates(
		ctx context.Context,
		input CalculateRatesInput,
	) ([]RateOption, error)
}

type LocationProvider interface {
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

type CreateOrderInput struct {
	OriginAreaID      int
	DestinationAreaID int
	UniqueOrderID     string
	CourierCode       string
	CourierService    string
	Weight            int

	ItemName  string
	ItemPrice int64
	ItemQty   int

	ShipperName     string
	ShipperPhone    string
	ShipperAddress  string
	ReceiverName    string
	ReceiverPhone   string
	ReceiverAddress string

	// ManualTrackingNumber is an optional pre-set
	// tracking number supplied by staff when the server
	// runs in manual logistics mode. Automated providers
	// (e.g. Komerce) ignore this field entirely.
	ManualTrackingNumber *string
}

type CreateOrderResult struct {
	KomerceOrderNo string
	TrackingNumber string
}

type TrackShipmentInput struct {
	Courier        string
	TrackingNumber string
	LastPhone      *string
}

type TrackingEvent struct {
	Status      string
	Description string
	Location    string
	Timestamp   time.Time
}

type LogisticsProvider interface {
	CreateOrder(
		ctx context.Context,
		input CreateOrderInput,
	) (*CreateOrderResult, error)

	CancelOrder(
		ctx context.Context,
		komerceOrderNo string,
	) error

	TrackShipment(
		ctx context.Context,
		input TrackShipmentInput,
	) ([]TrackingEvent, error)
}
