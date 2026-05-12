package bootstrap

import (
	"service-core/internal/common/logger"

	authDomain "service-core/internal/modules/auth/domain"
	authService "service-core/internal/modules/auth/infra/service"

	imgService "service-core/internal/shared/image"
	sGen "service-core/internal/shared/slug"

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
	shUC "service-core/internal/modules/shipment/usecase"
	sUC "service-core/internal/modules/shop/usecase"
	uUC "service-core/internal/modules/user/usecase"
)

type Container struct {
	Logger logger.Logger

	TokenService authDomain.TokenService
	Hasher       authDomain.PasswordHasher

	ImageVariantProvider imgService.VariantCreator
	ImageTransformer     imgService.ImageTransformer

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

	EstimateShippingOptions shUC.EstimateShippingOptionsUsecase
}

func NewContainer(cfg Config, infra *Infra) *Container {
	var (
		log = logger.NewZapLogger(cfg.App.Env)

		productRepo      = pRepoImpl.NewProductRepository(infra.DB)
		productImageRepo = pRepoImpl.NewProductImageRepository(infra.DB)
		inventoryRepo    = iRepoImpl.NewInventoryRepository(infra.DB)
		authRepo         = aRepoImpl.NewAuthRepository(infra.DB)
		cartRepo         = cRepoImpl.NewCartRepositoryImpl(infra.DB)
		// locationRepo = lRepoImpl.NewLocationRepositoryImpl(db)
		userRepo          = uRepoImpl.NewUserRepositoryImpl(infra.DB)
		addressRepo       = adRepoImpl.NewUserAddressRepositoryImpl(infra.DB)
		addressShopRepo   = adRepoImpl.NewShopAddressRepositoryImpl(infra.DB)
		paymentAccRepo    = payRepoImpl.NewPaymentAccountRepository(infra.DB)
		paymentMethodRepo = payRepoImpl.NewPaymentMethodRepository(infra.DB)
		shopRepo          = sRepoImpl.NewShopRepositoryImpl(infra.DB)
		courierRepo       = coRepoImpl.NewCourierRepositoryImpl(infra.DB)
		shopCourierRepo   = coRepoImpl.NewShopCourierRepositoryImpl(infra.DB)
	)

	var (
		tokenSvc = authService.NewJWTService(
			cfg.JWT.Secret,
			cfg.JWT.Exp,
		)

		hasher = authService.NewBcryptHasher()

		slugGen = sGen.NewGenerator()

		imageTransformer     = imgService.NewImageTransformer()
		imageVariantProvider = imgService.NewResolutionGenerator(imageTransformer)
	)

	return &Container{
		Logger: log,

		TokenService: tokenSvc,
		Hasher:       hasher,

		FindProducts: *pUC.NewFindProductsUsecase(productRepo, inventoryRepo),
		GetProduct:   *pUC.NewGetProductUsecase(productRepo, inventoryRepo),
		CreateProduct: *pUC.NewCreateProductUsecase(
			productRepo, productImageRepo, slugGen, imageVariantProvider, infra.StorageProvider,
		),
		CreateInventory: *iUC.NewCreateInventoryUsecase(inventoryRepo, productRepo, shopRepo),

		LoginAccount:    *aUC.NewLoginEmailUsecase(authRepo, hasher, tokenSvc),
		RegisterAccount: *aUC.NewRegisterUsecase(authRepo, hasher),
		GetAccount:      *aUC.NewGetAccountUsecase(authRepo),

		GetCart:    *cUC.NewGetCartUsecase(cartRepo, inventoryRepo, productRepo),
		AddItem:    *cUC.NewAddItemUsecase(cartRepo, inventoryRepo, productRepo),
		UpdateItem: *cUC.NewUpdateItemUsecase(cartRepo, inventoryRepo, productRepo),
		RemoveItem: *cUC.NewRemoveItemUsecase(cartRepo),

		ListLocations: *lUC.NewListLocationUsecase(infra.LocationRepository),

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

		EstimateShippingOptions: *shUC.NewEstimateShippingOptionsUsecase(infra.ShippingCostProvider),
	}
}
