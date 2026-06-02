package bootstrap

import (
	"net/http"

	database "service-core/internal/infra/db"
	storage "service-core/internal/infra/storage"

	supabaseStorage "service-core/internal/infra/storage/supabase"
	lService "service-core/internal/modules/location/infra/service"
	locationRepo "service-core/internal/modules/location/repository"
	sCostService "service-core/internal/modules/shipment/infra/service"
	shipmentRepo "service-core/internal/modules/shipment/repository"
)

type Dependency struct {
	DB                   *database.Connection
	StorageProvider      storage.Provider
	LocationRepository   locationRepo.LocationRepository
	ShippingCostProvider shipmentRepo.ShippingCostProvider
}

func NewDependency(cfg Config) (*Dependency, error) {
	db, err := database.NewConnection(cfg.DB)
	if err != nil {
		return nil, err
	}

	storageProvider, err := supabaseStorage.NewSupabaseProvider(
		cfg.Storage,
		cfg.Supabase,
		&http.Client{},
	)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Dependency{
		DB:              db,
		StorageProvider: storageProvider,
		LocationRepository: lService.NewRajaOngkirLocation(
			cfg.Shipping.DestinationKey,
			cfg.Shipping.DestinationURL,
			&http.Client{
				Timeout: cfg.Shipping.Timeout,
			},
		),
		ShippingCostProvider: sCostService.NewRajaOngkirCostEstimation(
			cfg.Shipping.CalculateKEY,
			cfg.Shipping.CalculateURL,
			&http.Client{
				Timeout: cfg.Shipping.Timeout,
			},
		),
	}, nil
}

func (i *Dependency) Close() {
	if i == nil || i.DB == nil {
		return
	}

	i.DB.Close()
}
