package bootstrap

import (
	"time"

	"service-core/internal/common/logger"
	"service-core/internal/shared/config"

	authDomain "service-core/internal/modules/auth/domain"
	authService "service-core/internal/modules/auth/infra/service"

	database "service-core/internal/infra/db"

	adRepoImpl "service-core/internal/modules/address/infra/persistence"
	aRepoImpl "service-core/internal/modules/auth/infra/persistence"
	cRepoImpl "service-core/internal/modules/cart/infra/persistence"
	lRepoImpl "service-core/internal/modules/location/infra/persistence"
	payRepoImpl "service-core/internal/modules/payment/infra/persistence"
	pRepoImpl "service-core/internal/modules/product/infra/persistence"
	uRepoImpl "service-core/internal/modules/user/infra/persistence"

	adUC "service-core/internal/modules/address/usecase"
	aUC "service-core/internal/modules/auth/usecase"
	cUC "service-core/internal/modules/cart/usecase"
	lUC "service-core/internal/modules/location/usecase"
	payUC "service-core/internal/modules/payment/usecase"
	pUC "service-core/internal/modules/product/usecase"
	uUC "service-core/internal/modules/user/usecase"
)

type Container struct {
	DB     *database.Connection
	Logger logger.Logger

	TokenService authDomain.TokenService
	Hasher       authDomain.PasswordHasher

	FindProducts  pUC.FindProductsUsecase
	GetProduct    pUC.GetProductUsecase
	CreateProduct pUC.CreateProductUsecase

	LoginAccount    aUC.LoginUsecase
	RegisterAccount aUC.RegisterUsecase
	GetAccount      aUC.GetAccountUsecase

	GetCart    cUC.GetCartUsecase
	AddItem    cUC.AddItemUsecase
	UpdateItem cUC.UpdateItemUsecase
	RemoveItem cUC.RemoveItemUsecase

	ListLocations lUC.ListLocationUsecase
	GetUser       uUC.GetUserUsecase
	GetAddress    adUC.GetAddressUsecase
	CreateAddress adUC.CreateAddressUsecase

	CreatePaymentAccount payUC.CreatePaymentAccount
	ListPaymentAccount   payUC.ListPaymentAccount
	CreatePaymentMethod  payUC.CreatePaymentMethod
	ListPaymentMethod    payUC.ListPaymentMethod
}

func NewContainer() *Container {
	log := logger.NewSlogLogger(config.MustGetEnv("APP_ENV"))

	cfg := database.LoadConfig()
	db := database.NewConnection(cfg)

	productRepo := pRepoImpl.NewProductRepository(db)
	authRepo := aRepoImpl.NewAuthRepository(db)
	cartRepo := cRepoImpl.NewCartRepositoryImpl(db)
	locationRepo := lRepoImpl.NewLocationRepositoryImpl(db)
	userRepo := uRepoImpl.NewUserRepositoryImpl(db)
	addressRepo := adRepoImpl.NewAddressRepositoryImpl(db)
	paymentAccRepo := payRepoImpl.NewPaymentAccountRepository(db)
	paymentMethodRepo := payRepoImpl.NewPaymentMethodRepository(db)

	tokenSvc := authService.NewJWTService(
		config.MustGetEnv("JWT_SECRET"),
		24*time.Minute,
	)

	hasher := authService.NewBcryptHasher()

	return &Container{
		DB:     db,
		Logger: log,

		TokenService: tokenSvc,
		Hasher:       hasher,

		FindProducts:  *pUC.NewFindProductsUsecase(productRepo),
		GetProduct:    *pUC.NewGetProductsUsecase(productRepo),
		CreateProduct: *pUC.NewCreateProductUsecase(productRepo),

		LoginAccount:    *aUC.NewLoginUsecase(authRepo, hasher, tokenSvc),
		RegisterAccount: *aUC.NewRegisterUsecase(authRepo, hasher),
		GetAccount:      *aUC.NewGetAccountUsecase(authRepo),

		GetCart:    *cUC.NewGetCartUsecase(cartRepo, productRepo),
		AddItem:    *cUC.NewAddItemUsecase(cartRepo, productRepo),
		UpdateItem: *cUC.NewUpdateItemUsecase(cartRepo, productRepo),
		RemoveItem: *cUC.NewRemoveItemUsecase(cartRepo),

		ListLocations: *lUC.NewListLocationUsecase(locationRepo),
		GetUser:       *uUC.NewGetUserUsecase(userRepo),
		GetAddress:    *adUC.NewGetAddressUsecase(addressRepo),
		CreateAddress: *adUC.NewCreateAddressUsecase(addressRepo),

		CreatePaymentAccount: *payUC.NewCreatePaymentAccount(paymentAccRepo, paymentMethodRepo),
		ListPaymentAccount:   *payUC.NewListPaymentAccount(paymentAccRepo),
		CreatePaymentMethod:  *payUC.NewCreatePaymentMethod(paymentMethodRepo),
		ListPaymentMethod:    *payUC.NewListPaymentMethod(paymentMethodRepo),
	}
}
