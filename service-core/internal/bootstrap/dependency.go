package bootstrap

import (
	"net/http"

	database "service-core/internal/infra/db"
	paymentgateway "service-core/internal/infra/payment-gateway"
	midtransGateway "service-core/internal/infra/payment-gateway/midtrans"
	"service-core/internal/infra/shipping"
	"service-core/internal/infra/shipping/rajaongkir"
	storage "service-core/internal/infra/storage"
	supabaseStorage "service-core/internal/infra/storage/supabase"
	lService "service-core/internal/modules/location/infra/service"
	locationRepo "service-core/internal/modules/location/repository"
	sCostService "service-core/internal/modules/shipment/infra/service"
	shipmentRepo "service-core/internal/modules/shipment/repository"
	transaction "service-core/internal/shared/transaction"
)

type Dependency struct {
	DB                   *database.Connection
	StorageProvider      storage.Provider
	LocationProvider     locationRepo.LocationRepository
	ShippingCostProvider shipmentRepo.ShippingCostProvider
	TransactionProvider  transaction.Transactor
	TransactionExecutor  transaction.Executor
	PaymentGateway       paymentgateway.Provider
	ShippingProvider     shipping.Provider
}

func NewDependency(cfg Config) (*Dependency, error) {
	storageProvider, err := supabaseStorage.
		NewSupabaseProvider(
			cfg.Storage,
			cfg.Supabase,
			&http.Client{},
		)
	if err != nil {
		return nil, err
	}

	gateway, err := midtransGateway.
		NewMidtransAPIProvider(cfg.MidTrans)
	if err != nil {
		return nil, err
	}

	shipping, err := rajaongkir.
		NewRajaOngkirProvider(cfg.RajaOngkir)
	if err != nil {
		return nil, err
	}

	db, err := database.
		NewConnection(cfg.DB)
	if err != nil {
		return nil, err
	}

	return &Dependency{
		DB:              db,
		StorageProvider: storageProvider,
		LocationProvider: lService.NewRajaOngkirLocation(
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
		TransactionProvider: database.NewPostgresTransactor(db.Pool),
		TransactionExecutor: db.Pool,
		PaymentGateway:      gateway,
		ShippingProvider:    shipping,
	}, nil
}

func (i *Dependency) Close() {
	if i == nil || i.DB == nil {
		return
	}

	i.DB.Close()
}
