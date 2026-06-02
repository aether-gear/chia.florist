package bootstrap

import (
	applogger "service-core/internal/common/logger"

	imgSvc "service-core/internal/shared/image"
	mailerSvc "service-core/internal/shared/mailer"
	otpSvc "service-core/internal/shared/otp"
	sGen "service-core/internal/shared/slug"

	addressPersistence "service-core/internal/modules/address/infra/persistence"
	authenPersistence "service-core/internal/modules/authentication/infra/persistence"
	cartPersistence "service-core/internal/modules/cart/infra/persistence"
	courierPersistence "service-core/internal/modules/courier/infra/persistence"
	inventoryPersistence "service-core/internal/modules/inventory/infra/persistence"
	merchantPersistence "service-core/internal/modules/merchant/infra/persistence"
	paymentPersistence "service-core/internal/modules/payment/infra/persistence"
	productPersistence "service-core/internal/modules/product/infra/persistence"
	shopPersistence "service-core/internal/modules/shop/infra/persistence"
	userPersistence "service-core/internal/modules/user/infra/persistence"

	authenSvc "service-core/internal/modules/authentication/infra/service"
	authenRepo "service-core/internal/modules/authentication/repository"
	authorSvc "service-core/internal/modules/authorization/infra/service"
	authorRepo "service-core/internal/modules/authorization/repository"

	addressUsecase "service-core/internal/modules/address/usecase"
	authenUsecase "service-core/internal/modules/authentication/usecase"
	cartUsecase "service-core/internal/modules/cart/usecase"
	courierUsecase "service-core/internal/modules/courier/usecase"
	inventoryUsecase "service-core/internal/modules/inventory/usecase"
	locationUsecase "service-core/internal/modules/location/usecase"
	paymentUsecase "service-core/internal/modules/payment/usecase"
	productUsecase "service-core/internal/modules/product/usecase"
	shipmentUsecase "service-core/internal/modules/shipment/usecase"
	shopUsecase "service-core/internal/modules/shop/usecase"
	userUsecase "service-core/internal/modules/user/usecase"
)

type Container struct {
	Logger             applogger.Logger
	CORSAllowedOrigins []string
	Authenticator      authenRepo.Authenticator
	Authorizer         authorRepo.Authorizer

	ImageVariantProvider imgSvc.VariantCreator
	ImageTransformer     imgSvc.ImageTransformer

	FindProducts     productUsecase.FindProductsUsecase
	GetProduct       productUsecase.GetProductUsecase
	CreateProduct    productUsecase.CreateProductUsecase
	AddProductImages productUsecase.AddProductImagesUsecase
	CreateInventory  inventoryUsecase.CreateInventoryUsecase

	LoginCustomer    authenUsecase.LoginCustomerUsecase
	RegisterCustomer authenUsecase.RegisterCustomerUsecase
	VerifyAccount    authenUsecase.VerifyAccountUsecase
	GetAccount       authenUsecase.GetAccountUsecase

	GetCart    cartUsecase.GetCartUsecase
	AddItem    cartUsecase.AddItemUsecase
	UpdateItem cartUsecase.UpdateItemUsecase
	RemoveItem cartUsecase.RemoveItemUsecase

	ListLocations locationUsecase.ListLocationUsecase

	GetUser userUsecase.GetUserUsecase

	ListUserAddresses addressUsecase.ListUserAddressUsecase
	CreateAddress     addressUsecase.CreateAddressUsecase

	GetShopAddress    addressUsecase.GetShopAddressUsecase
	ListShopAddresses addressUsecase.ListShopAddressesUsecase
	CreateShopAddress addressUsecase.CreateShopAddressUsecase

	GetShop    shopUsecase.GetShopUsecase
	CreateShop shopUsecase.CreateShopUsecase

	CreatePaymentAccount paymentUsecase.CreatePaymentAccountUsecase
	ListPaymentAccount   paymentUsecase.ListPaymentAccountUsecase
	CreatePaymentMethod  paymentUsecase.CreatePaymentMethodUsecase
	ListPaymentMethod    paymentUsecase.ListPaymentMethodUsecase

	ConfigureShopCourier courierUsecase.ConfigureShopCourierUsecase

	EstimateShippingOptions shipmentUsecase.EstimateShippingOptionsUsecase
}

func NewContainer(cfg Config, infra *Dependency) *Container {
	var (
		log = applogger.NewZapLogger(cfg.App.Env)
	)

	var (
		productRepo       = productPersistence.NewProductRepository(infra.DB)
		productImageRepo  = productPersistence.NewProductImageRepository(infra.DB)
		inventoryRepo     = inventoryPersistence.NewInventoryRepository(infra.DB)
		accountRepo       = authenPersistence.NewAccountRepository(infra.DB)
		challengeRepo     = authenPersistence.NewChallengeRepository(infra.DB)
		sessionRepo       = authenPersistence.NewSessionRepositoryImpl(infra.DB)
		refreshTokenRepo  = authenPersistence.NewRefreshTokenRepositoryImpl(infra.DB)
		cartRepo          = cartPersistence.NewCartRepositoryImpl(infra.DB)
		userRepo          = userPersistence.NewUserRepositoryImpl(infra.DB)
		addressRepo       = addressPersistence.NewUserAddressRepositoryImpl(infra.DB)
		addressShopRepo   = addressPersistence.NewShopAddressRepositoryImpl(infra.DB)
		paymentAccRepo    = paymentPersistence.NewPaymentAccountRepository(infra.DB)
		paymentMethodRepo = paymentPersistence.NewPaymentMethodRepository(infra.DB)
		shopRepo          = shopPersistence.NewShopRepositoryImpl(infra.DB)
		courierRepo       = courierPersistence.NewCourierRepositoryImpl(infra.DB)
		shopCourierRepo   = courierPersistence.NewShopCourierRepositoryImpl(infra.DB)
		merchantRepo      = merchantPersistence.NewMerchantRepositoryImpl(infra.DB)
	)

	var (
		tokenSvc    = authenSvc.NewJWTService(cfg.JWT.Secret)
		pwHasher    = authenSvc.NewBcryptHasher()
		tokenHasher = authenSvc.NewSHATokenHasher()
		authMidd    = authenSvc.NewJWTAuthenticator(tokenSvc, sessionRepo)

		actorSvc   = authorSvc.NewActorService(accountRepo, merchantRepo)
		authorMdwr = authorSvc.NewAuthorizer(actorSvc)
	)

	var (
		slugGen = sGen.NewGenerator()

		mailSender = mailerSvc.NewSMTPSender(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.Username,
			cfg.SMTP.Password,
			cfg.SMTP.From,
		)

		otpGen = otpSvc.NewNumericGenerator(6)

		imageTransformer     = imgSvc.NewImageTransformer()
		imageVariantProvider = imgSvc.NewResolutionGenerator(imageTransformer)
	)

	return &Container{
		Logger:             log,
		CORSAllowedOrigins: cfg.App.CORSAllowedOrigins,
		Authenticator:      authMidd,
		Authorizer:         authorMdwr,

		FindProducts: *productUsecase.NewFindProductsUsecase(
			productRepo, inventoryRepo, productImageRepo, infra.StorageProvider,
		),
		GetProduct: *productUsecase.NewGetProductUsecase(
			productRepo, inventoryRepo, productImageRepo, infra.StorageProvider,
		),
		CreateProduct: *productUsecase.NewCreateProductUsecase(productRepo, slugGen),
		AddProductImages: *productUsecase.NewAddProductImagesUsecase(
			productRepo, productImageRepo, slugGen, imageVariantProvider, infra.StorageProvider,
		),
		CreateInventory: *inventoryUsecase.NewCreateInventoryUsecase(inventoryRepo, productRepo, shopRepo),

		LoginCustomer: *authenUsecase.NewLoginCustomerUsecase(
			accountRepo, pwHasher, tokenHasher, tokenSvc, sessionRepo, refreshTokenRepo,
		),
		RegisterCustomer: *authenUsecase.NewRegisterCustomerUsecase(
			accountRepo, pwHasher, userRepo, challengeRepo, otpGen, mailSender,
		),
		VerifyAccount: *authenUsecase.NewVerifyAccountUsecase(
			accountRepo, pwHasher, tokenHasher, userRepo, challengeRepo, tokenSvc, sessionRepo, refreshTokenRepo,
		),
		GetAccount: *authenUsecase.NewGetAccountUsecase(accountRepo),

		GetCart: *cartUsecase.NewGetCartUsecase(
			cartRepo, inventoryRepo, productRepo, productImageRepo, infra.StorageProvider,
		),
		AddItem:    *cartUsecase.NewAddItemUsecase(cartRepo, inventoryRepo, productRepo),
		UpdateItem: *cartUsecase.NewUpdateItemUsecase(cartRepo, inventoryRepo, productRepo),
		RemoveItem: *cartUsecase.NewRemoveItemUsecase(cartRepo),

		ListLocations: *locationUsecase.NewListLocationUsecase(infra.LocationRepository),

		GetUser: *userUsecase.NewGetUserUsecase(userRepo),

		ListUserAddresses: *addressUsecase.NewListUserAddressUsecase(addressRepo),
		CreateAddress:     *addressUsecase.NewCreateAddressUsecase(addressRepo),

		GetShopAddress:    *addressUsecase.NewGetShopAddressUsecase(addressShopRepo),
		ListShopAddresses: *addressUsecase.NewListShopAddressesUsecase(addressShopRepo),
		CreateShopAddress: *addressUsecase.NewCreateShopAddressUsecase(addressShopRepo),

		GetShop:    *shopUsecase.NewGetShopUsecase(shopRepo),
		CreateShop: *shopUsecase.NewCreateShopUsecase(shopRepo, slugGen),

		CreatePaymentAccount: *paymentUsecase.NewCreatePaymentAccountUsecase(paymentAccRepo, paymentMethodRepo),
		ListPaymentAccount:   *paymentUsecase.NewListPaymentAccountUsecase(paymentAccRepo),
		CreatePaymentMethod:  *paymentUsecase.NewCreatePaymentMethodUsecase(paymentMethodRepo),
		ListPaymentMethod:    *paymentUsecase.NewListPaymentMethodUsecase(paymentMethodRepo),

		ConfigureShopCourier: *courierUsecase.NewConfigureShopCourierUsecase(courierRepo, shopCourierRepo, shopRepo),

		EstimateShippingOptions: *shipmentUsecase.NewEstimateShippingOptionsUsecase(infra.ShippingCostProvider),
	}
}
