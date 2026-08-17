package bootstrap

import (
	"fmt"
	"net/http"

	database "service-core/internal/infra/db"
	paymentgateway "service-core/internal/infra/payment-gateway"
	midtransGateway "service-core/internal/infra/payment-gateway/midtrans"
	"service-core/internal/infra/shipping"
	komerceProvider "service-core/internal/infra/shipping/komerce"
	manualShipProvider "service-core/internal/infra/shipping/manual"
	"service-core/internal/infra/shipping/rajaongkir"
	storage "service-core/internal/infra/storage"
	supabaseStorage "service-core/internal/infra/storage/supabase"
	transaction "service-core/internal/shared/transaction"
)

type Dependency struct {
	DB                  *database.Connection
	StorageProvider     storage.Provider
	TransactionProvider transaction.Transactor
	TransactionExecutor transaction.Executor
	PaymentGateway      paymentgateway.Provider
	LocationProvider    shipping.LocationProvider
	ShippingProvider    shipping.ShippingProvider
	LogisticsProvider   shipping.LogisticsProvider
}

func NewDependency(cfg Config) (*Dependency, error) {
	storage, err := supabaseStorage.NewSupabaseProvider(
		cfg.Storage,
		cfg.Supabase,
		&http.Client{},
	)
	if err != nil {
		return nil, err
	}

	gateway, err := midtransGateway.NewMidtransAPIProvider(cfg.MidTrans)
	if err != nil {
		return nil, err
	}

	shippingEstimator, err := rajaongkir.NewRajaOngkirProvider(cfg.RajaOngkir)
	if err != nil {
		return nil, err
	}

	location, err := rajaongkir.NewRajaOngkirProvider(cfg.RajaOngkir)
	if err != nil {
		return nil, err
	}

	var logistics shipping.LogisticsProvider
	switch cfg.Logistics.Provider {
	case "", "manual":
		logistics = manualShipProvider.NewManualShippingProvider()
	case "manual_tracked", "manual_external":
		komerce, err := komerceProvider.NewKomerceProvider(cfg.Komerce)
		if err != nil {
			return nil, err
		}
		logistics = manualShipProvider.NewManualTrackedShippingProvider(komerce)
	case "komerce":
		logistics, err = komerceProvider.NewKomerceProvider(cfg.Komerce)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported logistics provider: %s", cfg.Logistics.Provider)
	}

	db, err := database.NewConnection(cfg.DB)
	if err != nil {
		return nil, err
	}

	return &Dependency{
		DB:                  db,
		StorageProvider:     storage,
		TransactionProvider: database.NewPostgresTransactor(db.Pool),
		TransactionExecutor: db.Pool,
		PaymentGateway:      gateway,
		LocationProvider:    location,
		ShippingProvider:    shippingEstimator,
		LogisticsProvider:   logistics,
	}, nil
}

func (i *Dependency) Close() {
	if i == nil || i.DB == nil {
		return
	}

	i.DB.Close()
}
