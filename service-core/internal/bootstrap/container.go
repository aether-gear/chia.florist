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
	coRepoImpl "service-core/internal/modules/courier/infra/persistence"
	iRepoImpl "service-core/internal/modules/inventory/infra/persistence"
	payRepoImpl "service-core/internal/modules/payment/infra/persistence"
	pRepoImpl "service-core/internal/modules/product/infra/persistence"
	sRepoImpl "service-core/internal/modules/shop/infra/persistence"
	uRepoImpl "service-core/internal/modules/user/infra/persistence"

	adUC "service-core/internal/modules/address/usecase"
	aUC "service-core/internal/modules/auth/usecase"
	cUC "service-core/internal/modules/cart/usecase"
	coUC "service-core/internal/modules/courier/usecase"
	iUC "service-core/internal/modules/inventory/usecase"
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

	FindProducts    pUC.FindProductsUsecase
	GetProduct      pUC.GetProductUsecase
	CreateProduct   pUC.CreateProductUsecase
	CreateInventory iUC.CreateInventoryUsecase

	LoginAccount    aUC.LoginEmailUsecase
	RegisterAccount aUC.RegisterUsecase
	GetAccount      aUC.GetAccountUsecase

	GetCart    cUC.GetCartUsecase
	AddItem    cUC.AddItemUsecase
	UpdateItem cUC.UpdateItemUsecase
	RemoveItem cUC.RemoveItemUsecase

	ListLocations lUC.ListLocationUsecase
	GetUser       uUC.GetUserUsecase

	ListUserAddresses adUC.ListUserAddressUsecase
	CreateAddress     adUC.CreateAddressUsecase

	GetShopAddress    adUC.GetShopAddressUsecase
	ListShopAddresses adUC.ListShopAddressesUsecase
	CreateShopAddress adUC.CreateShopAddressUsecase

	GetShop    sUC.GetShopUsecase
	CreateShop sUC.CreateShopUsecase

	CreatePaymentAccount payUC.CreatePaymentAccountUsecase
	ListPaymentAccount   payUC.ListPaymentAccountUsecase
	CreatePaymentMethod  payUC.CreatePaymentMethodUsecase
	ListPaymentMethod    payUC.ListPaymentMethodUsecase

	ConfigureShopCourier coUC.ConfigureShopCourierUsecase
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

		productRepo   = pRepoImpl.NewProductRepository(db)
		inventoryRepo = iRepoImpl.NewInventoryRepository(db)
		authRepo      = aRepoImpl.NewAuthRepository(db)
		cartRepo      = cRepoImpl.NewCartRepositoryImpl(db)
		// locationRepo = lRepoImpl.NewLocationRepositoryImpl(db)
		userRepo          = uRepoImpl.NewUserRepositoryImpl(db)
		addressRepo       = adRepoImpl.NewUserAddressRepositoryImpl(db)
		addressShopRepo   = adRepoImpl.NewShopAddressRepositoryImpl(db)
		paymentAccRepo    = payRepoImpl.NewPaymentAccountRepository(db)
		paymentMethodRepo = payRepoImpl.NewPaymentMethodRepository(db)
		shopRepo          = sRepoImpl.NewShopRepositoryImpl(db)
		courierRepo       = coRepoImpl.NewCourierRepositoryImpl(db)
		shopCourierRepo   = coRepoImpl.NewShopCourierRepositoryImpl(db)
	)

	var (
		tokenSvc = authService.NewJWTService(
			jwtSecret,
			jwtExp,
		)

		hasher = authService.NewBcryptHasher()

		locationService = lService.NewRajaOngkirLocation(
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

		FindProducts:    *pUC.NewFindProductsUsecase(productRepo, inventoryRepo),
		GetProduct:      *pUC.NewGetProductUsecase(productRepo, inventoryRepo),
		CreateProduct:   *pUC.NewCreateProductUsecase(productRepo, slugGen),
		CreateInventory: *iUC.NewCreateInventoryUsecase(inventoryRepo, productRepo, shopRepo),

		LoginAccount:    *aUC.NewLoginEmailUsecase(authRepo, hasher, tokenSvc),
		RegisterAccount: *aUC.NewRegisterUsecase(authRepo, hasher),
		GetAccount:      *aUC.NewGetAccountUsecase(authRepo),

		GetCart:    *cUC.NewGetCartUsecase(cartRepo, inventoryRepo, productRepo),
		AddItem:    *cUC.NewAddItemUsecase(cartRepo, inventoryRepo, productRepo),
		UpdateItem: *cUC.NewUpdateItemUsecase(cartRepo, inventoryRepo, productRepo),
		RemoveItem: *cUC.NewRemoveItemUsecase(cartRepo),

		ListLocations: *lUC.NewListLocationUsecase(locationService),

		GetUser: *uUC.NewGetUserUsecase(userRepo),

		ListUserAddresses: *adUC.NewListUserAddressUsecase(addressRepo),
		CreateAddress:     *adUC.NewCreateAddressUsecase(addressRepo),

		GetShopAddress:    *adUC.NewGetShopAddressUsecase(addressShopRepo),
		ListShopAddresses: *adUC.NewListShopAddressesUsecase(addressShopRepo),
		CreateShopAddress: *adUC.NewCreateShopAddressUsecase(addressShopRepo),

		GetShop:    *sUC.NewGetShopUsecase(shopRepo),
		CreateShop: *sUC.NewCreateShopUsecase(shopRepo, slugGen),

		CreatePaymentAccount: *payUC.NewCreatePaymentAccountUsecase(paymentAccRepo, paymentMethodRepo),
		ListPaymentAccount:   *payUC.NewListPaymentAccountUsecase(paymentAccRepo),
		CreatePaymentMethod:  *payUC.NewCreatePaymentMethodUsecase(paymentMethodRepo),
		ListPaymentMethod:    *payUC.NewListPaymentMethodUsecase(paymentMethodRepo),

		ConfigureShopCourier: *coUC.NewConfigureShopCourierUsecase(courierRepo, shopCourierRepo, shopRepo),
	}
}
