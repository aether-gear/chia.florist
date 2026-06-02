package bootstrap

import (
	"service-core/internal/common/logger"
	appmiddleware "service-core/internal/common/middleware"

	imgService "service-core/internal/shared/image"
	mailer "service-core/internal/shared/mailer"
	otp "service-core/internal/shared/otp"
	sGen "service-core/internal/shared/slug"

	adRepoImpl "service-core/internal/modules/address/infra/persistence"
	aRepoImpl "service-core/internal/modules/authentication/infra/persistence"
	cRepoImpl "service-core/internal/modules/cart/infra/persistence"
	coRepoImpl "service-core/internal/modules/courier/infra/persistence"
	iRepoImpl "service-core/internal/modules/inventory/infra/persistence"
	mRepoImpl "service-core/internal/modules/merchant/infra/persistence"
	payRepoImpl "service-core/internal/modules/payment/infra/persistence"
	pRepoImpl "service-core/internal/modules/product/infra/persistence"
	sRepoImpl "service-core/internal/modules/shop/infra/persistence"
	uRepoImpl "service-core/internal/modules/user/infra/persistence"

	aService "service-core/internal/modules/authentication/infra/service"
	authorSvc "service-core/internal/modules/authorization/infra/service"
	authorRepo "service-core/internal/modules/authorization/repository"

	adUC "service-core/internal/modules/address/usecase"
	aUC "service-core/internal/modules/authentication/usecase"
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
	Logger             logger.Logger
	CORSAllowedOrigins []string
	AuthMiddleware     appmiddleware.Middleware
	Authorizer         authorRepo.Authorizer

	ImageVariantProvider imgService.VariantCreator
	ImageTransformer     imgService.ImageTransformer

	FindProducts     pUC.FindProductsUsecase
	GetProduct       pUC.GetProductUsecase
	CreateProduct    pUC.CreateProductUsecase
	AddProductImages pUC.AddProductImagesUsecase
	CreateInventory  iUC.CreateInventoryUsecase

	LoginAccount    aUC.LoginEmailUsecase
	RegisterAccount aUC.RegisterUsecase
	VerifyAccount   aUC.VerifyAccountUsecase
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
	)

	var (
		productRepo       = pRepoImpl.NewProductRepository(infra.DB)
		productImageRepo  = pRepoImpl.NewProductImageRepository(infra.DB)
		inventoryRepo     = iRepoImpl.NewInventoryRepository(infra.DB)
		authRepo          = aRepoImpl.NewAccountRepository(infra.DB)
		challengeRepo     = aRepoImpl.NewChallengeRepository(infra.DB)
		sessionRepo       = aRepoImpl.NewSessionRepositoryImpl(infra.DB)
		refreshTokenRepo  = aRepoImpl.NewRefreshTokenRepositoryImpl(infra.DB)
		cartRepo          = cRepoImpl.NewCartRepositoryImpl(infra.DB)
		userRepo          = uRepoImpl.NewUserRepositoryImpl(infra.DB)
		addressRepo       = adRepoImpl.NewUserAddressRepositoryImpl(infra.DB)
		addressShopRepo   = adRepoImpl.NewShopAddressRepositoryImpl(infra.DB)
		paymentAccRepo    = payRepoImpl.NewPaymentAccountRepository(infra.DB)
		paymentMethodRepo = payRepoImpl.NewPaymentMethodRepository(infra.DB)
		shopRepo          = sRepoImpl.NewShopRepositoryImpl(infra.DB)
		courierRepo       = coRepoImpl.NewCourierRepositoryImpl(infra.DB)
		shopCourierRepo   = coRepoImpl.NewShopCourierRepositoryImpl(infra.DB)
		merchantRepo      = mRepoImpl.NewMerchantRepositoryImpl(infra.DB)
	)

	var (
		tokenSvc    = aService.NewJWTService(cfg.JWT.Secret)
		pwHasher    = aService.NewBcryptHasher()
		tokenHasher = aService.NewSHATokenHasher()
		authMidd    = aService.NewAuthMiddleware(tokenSvc, sessionRepo)

		actorSvc   = authorSvc.NewActorService(authRepo, merchantRepo)
		authorMdwr = authorSvc.NewAuthorizerService(actorSvc)
	)

	var (
		slugGen = sGen.NewGenerator()

		mailSender = mailer.NewSMTPSender(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.Username,
			cfg.SMTP.Password,
			cfg.SMTP.From,
		)

		otpGen = otp.NewNumericGenerator(6)

		imageTransformer     = imgService.NewImageTransformer()
		imageVariantProvider = imgService.NewResolutionGenerator(imageTransformer)
	)

	return &Container{
		Logger:             log,
		CORSAllowedOrigins: cfg.App.CORSAllowedOrigins,
		AuthMiddleware:     authMidd.RequireAuth(),
		Authorizer:         authorMdwr,

		FindProducts: *pUC.NewFindProductsUsecase(
			productRepo, inventoryRepo, productImageRepo, infra.StorageProvider,
		),
		GetProduct: *pUC.NewGetProductUsecase(
			productRepo, inventoryRepo, productImageRepo, infra.StorageProvider,
		),
		CreateProduct: *pUC.NewCreateProductUsecase(productRepo, slugGen),
		AddProductImages: *pUC.NewAddProductImagesUsecase(
			productRepo, productImageRepo, slugGen, imageVariantProvider, infra.StorageProvider,
		),
		CreateInventory: *iUC.NewCreateInventoryUsecase(inventoryRepo, productRepo, shopRepo),

		LoginAccount: *aUC.NewLoginEmailUsecase(
			authRepo, pwHasher, tokenHasher, tokenSvc, sessionRepo, refreshTokenRepo,
		),
		RegisterAccount: *aUC.NewRegisterUsecase(
			authRepo, pwHasher, userRepo, challengeRepo, otpGen, mailSender,
		),
		VerifyAccount: *aUC.NewVerifyAccountUsecase(
			authRepo, pwHasher, tokenHasher, userRepo, challengeRepo, tokenSvc, sessionRepo, refreshTokenRepo,
		),
		GetAccount: *aUC.NewGetAccountUsecase(authRepo),

		GetCart: *cUC.NewGetCartUsecase(
			cartRepo, inventoryRepo, productRepo, productImageRepo, infra.StorageProvider,
		),
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
