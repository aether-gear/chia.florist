package bootstrap

import (
	"net/http"

	database "service-core/internal/infra/db"
	paymentgateway "service-core/internal/infra/payment-gateway"
	midtransGateway "service-core/internal/infra/payment-gateway/midtrans"
	"service-core/internal/infra/shipping"
	komerceProvider "service-core/internal/infra/shipping/komerce"
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
	ShippingProvider    shipping.Provider
	LogisticsProvider   shipping.LogisticsProvider
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

	logistics, err := komerceProvider.
		NewKomerceProvider(cfg.Komerce)
	if err != nil {
		return nil, err
	}

	db, err := database.
		NewConnection(cfg.DB)
	if err != nil {
		return nil, err
	}

	return &Dependency{
		DB:                  db,
		StorageProvider:     storageProvider,
		TransactionProvider: database.NewPostgresTransactor(db.Pool),
		TransactionExecutor: db.Pool,
		PaymentGateway:      gateway,
		ShippingProvider:    shipping,
		LogisticsProvider:   logistics,
	}, nil
}

func (i *Dependency) Close() {
	if i == nil || i.DB == nil {
		return
	}

	i.DB.Close()
}
