package bootstrap

import (
	"net/http"

	"service-core/internal/common/logger"

	authDomain "service-core/internal/modules/auth/domain"
	authService "service-core/internal/modules/auth/infra/service"
	lService "service-core/internal/modules/location/infra/service"
	sGen "service-core/internal/shared/slug"

	database "service-core/internal/infra/db"

	adRepoImpl "service-core/internal/modules/address/infra/persistence"
	aRepoImpl "service-core/internal/modules/auth/infra/persistence"
	cRepoImpl "service-core/internal/modules/cart/infra/persistence"

	// lRepoImpl "service-core/internal/modules/location/infra/persistence"
	payRepoImpl "service-core/internal/modules/payment/infra/persistence"
	pRepoImpl "service-core/internal/modules/product/infra/persistence"
	sRepoImpl "service-core/internal/modules/shop/infra/persistence"
	uRepoImpl "service-core/internal/modules/user/infra/persistence"

	adUC "service-core/internal/modules/address/usecase"
	aUC "service-core/internal/modules/auth/usecase"
	cUC "service-core/internal/modules/cart/usecase"
	lUC "service-core/internal/modules/location/usecase"
	payUC "service-core/internal/modules/payment/usecase"
	pUC "service-core/internal/modules/product/usecase"
	sUC "service-core/internal/modules/shop/usecase"
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

	GetShopAddress    adUC.GetShopAddressUsecase
	GetShopAddresses  adUC.FindShopAddressUsecase
	CreateShopAddress adUC.CreateShopAddressUsecase

	CreateShop sUC.CreateShopUsecase

	CreatePaymentAccount payUC.CreatePaymentAccount
	ListPaymentAccount   payUC.ListPaymentAccount
	CreatePaymentMethod  payUC.CreatePaymentMethod
	ListPaymentMethod    payUC.ListPaymentMethod
}

func NewContainer() *Container {
	var (
		config = LoadConfig()

		dbCfg                 = config.DB
		app                   = config.App.Env
		komerceShipping       = config.Shipping.BaseURL
		komerceDestinationURL = config.Shipping.DestinationURL
		komerceTimeout        = config.Shipping.Timeout
		jwtSecret             = config.JWT.Secret
		jwtExp                = config.JWT.Exp
	)

	var (
		db  = database.NewConnection(dbCfg)
		log = logger.NewZapLogger(app)

		productRepo = pRepoImpl.NewProductRepository(db)
		authRepo    = aRepoImpl.NewAuthRepository(db)
		cartRepo    = cRepoImpl.NewCartRepositoryImpl(db)
		// locationRepo = lRepoImpl.NewLocationRepositoryImpl(db)
		userRepo          = uRepoImpl.NewUserRepositoryImpl(db)
		addressRepo       = adRepoImpl.NewUserAddressRepositoryImpl(db)
		addressShopRepo   = adRepoImpl.NewShopAddressRepositoryImpl(db)
		paymentAccRepo    = payRepoImpl.NewPaymentAccountRepository(db)
		paymentMethodRepo = payRepoImpl.NewPaymentMethodRepository(db)
		shopRepo          = sRepoImpl.NewShopRepositoryImpl(db)
	)

	var (
		tokenSvc = authService.NewJWTService(
			jwtSecret,
			jwtExp,
		)

		hasher = authService.NewBcryptHasher()

		locationService = lService.NewRajaOngkirService(
			komerceShipping,
			komerceDestinationURL,
			&http.Client{
				Timeout: komerceTimeout,
			},
		)

		slugGen = sGen.NewGenerator()
	)

	return &Container{
		DB:     db,
		Logger: log,

		TokenService: tokenSvc,
		Hasher:       hasher,

		FindProducts:  *pUC.NewFindProductsUsecase(productRepo),
		GetProduct:    *pUC.NewGetProductsUsecase(productRepo),
		CreateProduct: *pUC.NewCreateProductUsecase(productRepo, slugGen),

		LoginAccount:    *aUC.NewLoginUsecase(authRepo, hasher, tokenSvc),
		RegisterAccount: *aUC.NewRegisterUsecase(authRepo, hasher),
		GetAccount:      *aUC.NewGetAccountUsecase(authRepo),

		GetCart:    *cUC.NewGetCartUsecase(cartRepo, productRepo),
		AddItem:    *cUC.NewAddItemUsecase(cartRepo, productRepo),
		UpdateItem: *cUC.NewUpdateItemUsecase(cartRepo, productRepo),
		RemoveItem: *cUC.NewRemoveItemUsecase(cartRepo),

		ListLocations: *lUC.NewListLocationUsecase(locationService),

		GetUser: *uUC.NewGetUserUsecase(userRepo),

		GetAddress:    *adUC.NewGetAddressUsecase(addressRepo),
		CreateAddress: *adUC.NewCreateAddressUsecase(addressRepo),

		GetShopAddress:    *adUC.NewGetShopAddressUsecase(addressShopRepo),
		GetShopAddresses:  *adUC.NewFindShopAddressUsecase(addressShopRepo),
		CreateShopAddress: *adUC.NewCreateShopAddressUsecase(addressShopRepo),

		CreateShop: *sUC.NewCreateShopUsecase(shopRepo, slugGen),

		CreatePaymentAccount: *payUC.NewCreatePaymentAccount(paymentAccRepo, paymentMethodRepo),
		ListPaymentAccount:   *payUC.NewListPaymentAccount(paymentAccRepo),
		CreatePaymentMethod:  *payUC.NewCreatePaymentMethod(paymentMethodRepo),
		ListPaymentMethod:    *payUC.NewListPaymentMethod(paymentMethodRepo),
	}
}
