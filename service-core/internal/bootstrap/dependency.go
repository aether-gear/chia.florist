package bootstrap

import (
	"net/http"

	database "service-core/internal/infra/db"
	paymentgateway "service-core/internal/infra/payment-gateway"
	midtransGateway "service-core/internal/infra/payment-gateway/midtrans"
	storage "service-core/internal/infra/storage"
	transaction "service-core/internal/shared/transaction"

	supabaseStorage "service-core/internal/infra/storage/supabase"
	lService "service-core/internal/modules/location/infra/service"
	locationRepo "service-core/internal/modules/location/repository"
	sCostService "service-core/internal/modules/shipment/infra/service"
	shipmentRepo "service-core/internal/modules/shipment/repository"
)

type Dependency struct {
	DB                   *database.Connection
	StorageProvider      storage.Provider
	LocationProvider     locationRepo.LocationRepository
	ShippingCostProvider shipmentRepo.ShippingCostProvider
	TransactionProvider  transaction.Transactor
	TransactionExecutor  transaction.Executor
	PaymentGateway       paymentgateway.Provider
}

func NewDependency(cfg Config) (*Dependency, error) {
	db, err := database.NewConnection(cfg.DB)
	if err != nil {
		return nil, err
	}

	storageProvider, err :=
		supabaseStorage.NewSupabaseProvider(
			cfg.Storage,
			cfg.Supabase,
			&http.Client{},
		)
	if err != nil {
		db.Close()
		return nil, err
	}

	gateway, err :=
		midtransGateway.NewMidtransProvider(
			cfg.MidTrans,
		)
	if err != nil {
		db.Close()
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
	}, nil
}

func (i *Dependency) Close() {
	if i == nil || i.DB == nil {
		return
	}

	i.DB.Close()
}
