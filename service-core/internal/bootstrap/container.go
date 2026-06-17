package bootstrap

import (
	applogger "service-core/internal/common/logger"

	imgSvc "service-core/internal/shared/image"
	mailerSvc "service-core/internal/shared/mailer"
	otpSvc "service-core/internal/shared/otp"
	sGen "service-core/internal/shared/slug"
	"service-core/internal/shared/transaction"

	addressPersistence "service-core/internal/modules/address/infra/persistence"
	authenPersistence "service-core/internal/modules/authentication/infra/persistence"
	authorPersistence "service-core/internal/modules/authorization/infra/persistence"
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
	merchantUsecase "service-core/internal/modules/merchant/usecase"
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
	DBExecutor         transaction.Executor

	FindProducts     productUsecase.FindProductsUsecase
	GetProduct       productUsecase.GetProductUsecase
	CreateProduct    productUsecase.CreateProductUsecase
	AddProductImages productUsecase.AddProductImagesUsecase
	CreateInventory  inventoryUsecase.CreateInventoryUsecase

	Me               authenUsecase.MeUsecase
	LoginCustomer    authenUsecase.LoginCustomerUsecase
	LoginMerchant    authenUsecase.LoginMerchantUsecase
	RegisterCustomer authenUsecase.RegisterCustomerUsecase
	VerifyAccount    authenUsecase.VerifyAccountUsecase
	GetAccount       authenUsecase.GetAccountUsecase
	Logout           authenUsecase.LogoutUsecase

	CreateMerchant     merchantUsecase.CreateMerchantUsecase
	AddMerchantAccount merchantUsecase.AddMerchantAccountUsecase

	GetCart    cartUsecase.GetCartUsecase
	AddItem    cartUsecase.AddItemUsecase
	UpdateItem cartUsecase.UpdateItemUsecase
	RemoveItem cartUsecase.RemoveItemUsecase
	Checkout   cartUsecase.CheckoutUsecase

	ListLocations locationUsecase.ListLocationUsecase

	GetUser userUsecase.GetUserUsecase

	ListUserAddresses addressUsecase.ListUserAddressUsecase
	CreateUserAddress addressUsecase.SaveUserAddressUsecase
	DeleteUserAddress addressUsecase.DeleteUserAddressUsecase

	ListShopAddresses addressUsecase.ListShopAddressesUsecase
	SaveShopAddress   addressUsecase.CreateShopAddressUsecase

	FindShops  shopUsecase.FindShopsUsecase
	GetShop    shopUsecase.GetShopUsecase
	CreateShop shopUsecase.CreateShopUsecase

	CreatePaymentAccount paymentUsecase.CreatePaymentAccountUsecase
	ListPaymentAccount   paymentUsecase.ListPaymentAccountUsecase
	CreatePaymentMethod  paymentUsecase.CreatePaymentMethodUsecase
	ListPaymentMethod    paymentUsecase.ListPaymentMethodUsecase

	ListAllCouriers      courierUsecase.ListCouriersUsecase
	ConfigureShopCourier courierUsecase.ConfigureShopCourierUsecase

	EstimateShippingOptions shipmentUsecase.EstimateShippingOptionsUsecase
}

func NewContainer(cfg Config,
	infra *Dependency) *Container {
	var (
		log = applogger.NewZapLogger(cfg.App.Env)
	)

	var (
		productRepo       = productPersistence.NewProductRepository()
		productImageRepo  = productPersistence.NewProductImageRepository()
		inventoryRepo     = inventoryPersistence.NewInventoryRepository()
		accountRepo       = authenPersistence.NewAccountRepository()
		challengeRepo     = authenPersistence.NewChallengeRepository()
		sessionRepo       = authenPersistence.NewSessionRepositoryImpl()
		refreshTokenRepo  = authenPersistence.NewRefreshTokenRepositoryImpl()
		cartRepo          = cartPersistence.NewCartRepositoryImpl()
		userRepo          = userPersistence.NewUserRepositoryImpl()
		addressRepo       = addressPersistence.NewUserAddressRepositoryImpl()
		addressShopRepo   = addressPersistence.NewShopAddressRepositoryImpl()
		paymentAccRepo    = paymentPersistence.NewPaymentAccountRepository()
		paymentMethodRepo = paymentPersistence.NewPaymentMethodRepository()
		shopRepo          = shopPersistence.NewShopRepositoryImpl()
		courierRepo       = courierPersistence.NewCourierRepositoryImpl()
		shopCourierRepo   = courierPersistence.NewShopCourierRepositoryImpl()
		merchantRepo      = merchantPersistence.NewMerchantRepositoryImpl()
		membershipRepo    = authorPersistence.NewMerchantMembershipRepositoryImpl()
		roleRepo          = authorPersistence.NewRoleRepositoryImpl()
	)

	var (
		tokenSvc    = authenSvc.NewJWTService(cfg.JWT.Secret)
		pwHasher    = authenSvc.NewBcryptHasher()
		tokenHasher = authenSvc.NewSHATokenHasher()
		authMidd    = authenSvc.NewJWTAuthenticator(
			tokenSvc,
			sessionRepo,
		)

		actorSvc = authorSvc.NewActorService(
			accountRepo,
			merchantRepo,
			membershipRepo,
		)
		authorMdwr = authorSvc.NewAuthorizer(
			actorSvc,
		)
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
		DBExecutor:         infra.TransactionExecutor,

		FindProducts: *productUsecase.
			NewFindProductsUsecase(
				productRepo,
				inventoryRepo,
				productImageRepo,
				shopRepo,
				infra.StorageProvider,
				infra.TransactionExecutor,
			),
		GetProduct: *productUsecase.
			NewGetProductUsecase(
				productRepo,
				inventoryRepo,
				productImageRepo,
				infra.StorageProvider,
				infra.TransactionExecutor,
			),
		CreateProduct: *productUsecase.
			NewCreateProductUsecase(
				productRepo,
				slugGen,
				infra.TransactionExecutor,
			),
		AddProductImages: *productUsecase.
			NewAddProductImagesUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				productRepo,
				productImageRepo,
				slugGen,
				imageVariantProvider,
				infra.StorageProvider,
			),
		CreateInventory: *inventoryUsecase.
			NewCreateInventoryUsecase(inventoryRepo,
				productRepo,
				shopRepo,
				infra.TransactionExecutor,
			),

		Me: *authenUsecase.NewMeUsecase(
			infra.TransactionExecutor,
			accountRepo,
			actorSvc,
		),
		Logout: *authenUsecase.NewLogoutUsecase(
			infra.TransactionProvider,
			refreshTokenRepo,
			sessionRepo,
		),
		LoginCustomer: *authenUsecase.
			NewLoginCustomerUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				tokenHasher,
				tokenSvc,
				sessionRepo,
				refreshTokenRepo,
			),
		LoginMerchant: *authenUsecase.
			NewLoginMerchantUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				tokenHasher,
				tokenSvc,
				sessionRepo,
				refreshTokenRepo,
				merchantRepo,
				membershipRepo,
			),

		CreateMerchant: *merchantUsecase.
			NewCreateMerchantUsecase(
				merchantRepo,
				infra.TransactionExecutor,
			),
		AddMerchantAccount: *merchantUsecase.
			NewAddMerchantAccountUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				userRepo,
				membershipRepo,
				roleRepo,
			),

		RegisterCustomer: *authenUsecase.
			NewRegisterCustomerUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				userRepo,
				challengeRepo,
				otpGen,
				mailSender,
			),
		VerifyAccount: *authenUsecase.
			NewVerifyAccountUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				tokenHasher,
				userRepo,
				challengeRepo,
				tokenSvc,
				sessionRepo,
				refreshTokenRepo,
			),
		GetAccount: *authenUsecase.
			NewGetAccountUsecase(
				accountRepo,
				infra.TransactionExecutor,
			),

		GetCart: *cartUsecase.
			NewGetCartUsecase(
				cartRepo,
				inventoryRepo,
				productRepo,
				productImageRepo,
				infra.StorageProvider,
				infra.TransactionExecutor,
			),
		AddItem: *cartUsecase.
			NewAddItemUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				cartRepo,
				inventoryRepo,
				productRepo,
			),
		UpdateItem: *cartUsecase.
			NewUpdateItemUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				cartRepo,
				inventoryRepo,
				productRepo,
			),
		RemoveItem: *cartUsecase.
			NewRemoveItemUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				cartRepo,
			),
		Checkout: *cartUsecase.
			NewCheckoutUsecase(
				infra.TransactionExecutor,
				addressRepo,
				shopCourierRepo,
				inventoryRepo,
				productRepo,
				infra.ShippingCostProvider,
				addressShopRepo,
			),

		ListLocations: *locationUsecase.
			NewListLocationUsecase(
				infra.LocationProvider,
				infra.TransactionExecutor,
			),

		GetUser: *userUsecase.
			NewGetUserUsecase(
				userRepo,
				infra.TransactionExecutor,
			),

		ListUserAddresses: *addressUsecase.
			NewListUserAddressUsecase(
				addressRepo,
				infra.TransactionExecutor,
			),
		CreateUserAddress: *addressUsecase.
			NewSaveUserAddressUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				addressRepo,
			),
		DeleteUserAddress: *addressUsecase.
			NewDeleteUserAddressUsecase(
				infra.TransactionExecutor,
				addressRepo,
			),

		ListShopAddresses: *addressUsecase.
			NewListShopAddressesUsecase(
				addressShopRepo,
				infra.TransactionExecutor,
			),
		SaveShopAddress: *addressUsecase.
			NewCreateShopAddressUsecase(
				addressShopRepo,
				infra.TransactionExecutor,
			),

		FindShops: *shopUsecase.
			NewFindShopsUsecase(
				infra.TransactionExecutor,
				shopRepo,
			),
		GetShop: *shopUsecase.
			NewGetShopUsecase(
				shopRepo,
				infra.TransactionExecutor,
			),
		CreateShop: *shopUsecase.
			NewCreateShopUsecase(
				shopRepo,
				slugGen,
				infra.TransactionExecutor,
			),

		CreatePaymentAccount: *paymentUsecase.
			NewCreatePaymentAccountUsecase(
				paymentAccRepo,
				paymentMethodRepo,
				infra.TransactionExecutor,
			),
		ListPaymentAccount: *paymentUsecase.
			NewListPaymentAccountUsecase(
				paymentAccRepo,
				infra.TransactionExecutor,
			),
		CreatePaymentMethod: *paymentUsecase.
			NewCreatePaymentMethodUsecase(
				paymentMethodRepo,
				infra.TransactionExecutor,
			),
		ListPaymentMethod: *paymentUsecase.
			NewListPaymentMethodUsecase(
				paymentMethodRepo,
				infra.TransactionExecutor,
			),

		ListAllCouriers: *courierUsecase.NewListCouriersUsecase(
			infra.TransactionExecutor,
			courierRepo,
		),
		ConfigureShopCourier: *courierUsecase.
			NewConfigureShopCourierUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				courierRepo,
				shopCourierRepo,
				shopRepo,
			),

		EstimateShippingOptions: *shipmentUsecase.
			NewEstimateShippingOptionsUsecase(
				infra.ShippingCostProvider,
				infra.TransactionExecutor,
			),
	}
}
