package bootstrap

import (
	"time"

	applimiter "service-core/internal/common/limiter"
	applogger "service-core/internal/common/logger"

	appconfig "service-core/internal/shared/config"
	imgSvc "service-core/internal/shared/image"
	mailerSvc "service-core/internal/shared/mailer"
	otpSvc "service-core/internal/shared/otp"
	sGen "service-core/internal/shared/slug"
	"service-core/internal/shared/transaction"

	auditInfra "service-core/internal/modules/audit/infra"
	auditPersistence "service-core/internal/modules/audit/infra/persistence"
	auditUsecase "service-core/internal/modules/audit/usecase"

	secPolicyPersistence "service-core/internal/modules/security_policy/infra/persistence"
	secPolicyUsecase "service-core/internal/modules/security_policy/usecase"

	addressPersistence "service-core/internal/modules/address/infra/persistence"
	authenPersistence "service-core/internal/modules/authentication/infra/persistence"
	authorPersistence "service-core/internal/modules/authorization/infra/persistence"
	cartPersistence "service-core/internal/modules/cart/infra/persistence"
	courierPersistence "service-core/internal/modules/courier/infra/persistence"
	customerPersistence "service-core/internal/modules/customer/infra/persistence"
	inventoryPersistence "service-core/internal/modules/inventory/infra/persistence"
	orderPersistence "service-core/internal/modules/order/infra/persistence"
	paymentPersistence "service-core/internal/modules/payment/infra/persistence"
	productPersistence "service-core/internal/modules/product/infra/persistence"
	shipmentPersistence "service-core/internal/modules/shipment/infra/persistence"
	shopPersistence "service-core/internal/modules/shop/infra/persistence"
	staffPersistence "service-core/internal/modules/staff/infra/persistence"
	userPersistence "service-core/internal/modules/user/infra/persistence"

	authenSvc "service-core/internal/modules/authentication/infra/service"
	authenRepo "service-core/internal/modules/authentication/repository"
	customerSvc "service-core/internal/modules/customer/infra/service"
	authorSvc "service-core/internal/modules/authorization/infra/service"
	authorRepo "service-core/internal/modules/authorization/repository"
	orderSvc "service-core/internal/modules/order/infra/service"

	addressUsecase "service-core/internal/modules/address/usecase"
	authenUsecase "service-core/internal/modules/authentication/usecase"
	cartUsecase "service-core/internal/modules/cart/usecase"
	courierUsecase "service-core/internal/modules/courier/usecase"
	customerUsecase "service-core/internal/modules/customer/usecase"
	inventoryUsecase "service-core/internal/modules/inventory/usecase"
	locationUsecase "service-core/internal/modules/location/usecase"
	orderUsecase "service-core/internal/modules/order/usecase"
	paymentUsecase "service-core/internal/modules/payment/usecase"
	productUsecase "service-core/internal/modules/product/usecase"
	shipmentUsecase "service-core/internal/modules/shipment/usecase"
	shopUsecase "service-core/internal/modules/shop/usecase"
	staffUsecase "service-core/internal/modules/staff/usecase"
	userUsecase "service-core/internal/modules/user/usecase"
)

type Container struct {
	Logger             applogger.Logger
	AuditLogger        applogger.AuditLogger
	CORSAllowedOrigins []string
	Authenticator      authenRepo.Authenticator
	Authorizer         authorRepo.Authorizer
	DBExecutor         transaction.Executor
	DBTransactor       transaction.Transactor
	GoogleOAuth        appconfig.GoogleOAuthConfig

	FindProducts     productUsecase.FindProductsUsecase
	GetProduct       productUsecase.GetProductUsecase
	CreateProduct    productUsecase.CreateProductUsecase
	AddProductImages productUsecase.AddProductImagesUsecase
	CreateInventory  inventoryUsecase.CreateInventoryUsecase

	Me                   authenUsecase.MeUsecase
	LoginCustomer        authenUsecase.LoginCustomerUsecase
	LoginStaff           authenUsecase.LoginStaffUsecase
	RegisterCustomer     authenUsecase.RegisterCustomerUsecase
	VerifyAccount        authenUsecase.VerifyAccountUsecase
	GetAccount           authenUsecase.GetAccountUsecase
	Logout               authenUsecase.LogoutUsecase
	AuthenticateOAuth    authenUsecase.AuthenticateOAuthUsecase
	RequestPasswordReset authenUsecase.RequestPasswordResetUsecase
	VerifyPasswordReset  authenUsecase.VerifyPasswordResetUsecase
	ResetPassword         authenUsecase.ResetPasswordUsecase
	DeleteCustomerAccount customerUsecase.DeleteCustomerAccountUsecase

	FindStaff       staffUsecase.FindStaffUsecase
	CreateStaff     staffUsecase.CreateStaffUsecase
	AddStaffAccount staffUsecase.AddStaffAccountUsecase

	GetCart    cartUsecase.GetCartUsecase
	AddItem    cartUsecase.AddItemUsecase
	UpdateItem cartUsecase.UpdateItemUsecase
	RemoveItem cartUsecase.RemoveItemUsecase
	Checkout   cartUsecase.CheckoutUsecase

	ListLocations locationUsecase.ListLocationUsecase

	GetUser              userUsecase.GetUserUsecase
	GetCurrentProfile    userUsecase.GetCurrentProfileUsecase
	UpdateCurrentProfile userUsecase.UpdateCurrentProfileUsecase

	FindCustomers customerUsecase.FindCustomersUsecase

	ListUserAddresses addressUsecase.ListCustomerAddressesUsecase
	CreateUserAddress addressUsecase.SaveCustomerAddressUsecase
	DeleteUserAddress addressUsecase.DeleteCustomerAddressUsecase

	ListShopAddresses addressUsecase.ListShopAddressesUsecase
	SaveShopAddress   addressUsecase.CreateShopAddressUsecase

	FindShops shopUsecase.FindShopsUsecase
	GetShop   shopUsecase.GetShopUsecase
	SaveShop  shopUsecase.SaveShopUsecase

	GetShopAddresses shopUsecase.GetShopAddressesUsecase
	GetShopCouriers  shopUsecase.GetShopCouriersUsecase
	GetShopProducts  shopUsecase.GetShopProductsUsecase

	CreatePaymentAccount  paymentUsecase.CreatePaymentAccountUsecase
	ListPaymentAccount    paymentUsecase.ListPaymentAccountUsecase
	CreatePaymentMethod   paymentUsecase.CreatePaymentMethodUsecase
	ListPaymentMethod     paymentUsecase.ListPaymentMethodUsecase
	ProcessPaymentWebhook paymentUsecase.ProcessPaymentWebhookUsecase
	ProcessManualPayment  paymentUsecase.ProcessManualPaymentUsecase

	ListAllCouriers      courierUsecase.ListCouriersUsecase
	ConfigureShopCourier courierUsecase.ConfigureShopCourierUsecase

	EstimateShippingOptions shipmentUsecase.EstimateShippingOptionsUsecase

	CreateOrder       orderUsecase.CreateOrderUsecase
	FindOrders        orderUsecase.FindOrdersUsecase
	GetOrder          orderUsecase.GetOrderUsecase
	UpdateOrderStatus orderUsecase.UpdateOrderStatusUsecase

	FindAuditLogs auditUsecase.FindAuditLogsUsecase
	GetAuditLog   auditUsecase.GetAuditLogUsecase
	DeleteAuditLogs auditUsecase.DeleteAuditLogsUsecase

	WAFAutoBanEnabled bool
	Limiter           applimiter.Limiter
	ListRules         secPolicyUsecase.ListRulesUsecase
	CreateRule        secPolicyUsecase.CreateRuleUsecase
	ToggleRule        secPolicyUsecase.ToggleRuleUsecase
	UpdateRule        secPolicyUsecase.UpdateRuleUsecase
	DeleteRule        secPolicyUsecase.DeleteRuleUsecase
	GetIPConfig       secPolicyUsecase.GetIPConfigUsecase
	UpdateIPAction    secPolicyUsecase.UpdateIPActionUsecase
	GetFilters        secPolicyUsecase.GetFiltersUsecase
	UpdateFilter      secPolicyUsecase.UpdateFilterUsecase
	InspectPayload    secPolicyUsecase.InspectPayloadUsecase
}

func NewContainer(cfg Config,
	infra *Dependency) *Container {
	var (
		log = applogger.NewZapLogger(cfg.App.Env)

		auditLogRepo = auditPersistence.NewAuditLogRepository()
		auditLogger  = auditInfra.NewDBAuditLogger(
			auditLogRepo,
			log,
			infra.TransactionExecutor,
		)
	)

	var (
		productRepo            = productPersistence.NewProductRepository()
		productImageRepo       = productPersistence.NewProductImageRepository()
		inventoryRepo          = inventoryPersistence.NewInventoryRepository()
		secPolicyRepo          = secPolicyPersistence.NewSecurityPolicyRepository()
		accountRepo            = authenPersistence.NewAccountRepository()
		challengeRepo          = authenPersistence.NewChallengeRepository()
		oauthRepo              = authenPersistence.NewOAuthConnectionRepository()
		sessionRepo            = authenPersistence.NewSessionRepositoryImpl()
		refreshTokenRepo       = authenPersistence.NewRefreshTokenRepositoryImpl()
		cartRepo               = cartPersistence.NewCartRepositoryImpl()
		userRepo               = userPersistence.NewUserRepositoryImpl()
		addressRepo            = addressPersistence.NewCustomerAddressRepositoryImpl()
		addressShopRepo        = addressPersistence.NewShopAddressRepositoryImpl()
		paymentRepo            = paymentPersistence.NewPaymentRepositoryImpl()
		paymentAccRepo         = paymentPersistence.NewPaymentAccountRepository()
		paymentMethodRepo      = paymentPersistence.NewPaymentMethodRepository()
		paymentEventRepo       = paymentPersistence.NewPaymentEventRepositoryImpl()
		paymentInstructionRepo = paymentPersistence.NewPaymentInstructionRepositoryImpl()
		shopRepo               = shopPersistence.NewShopRepositoryImpl()
		courierRepo            = courierPersistence.NewCourierRepositoryImpl()
		shopCourierRepo        = courierPersistence.NewShopCourierRepositoryImpl()
		staffRepo              = staffPersistence.NewStaffRepositoryImpl()
		customerRepo           = customerPersistence.NewCustomerRepositoryImpl()
		membershipRepo         = authorPersistence.NewStaffMembershipRepositoryImpl()
		roleRepo               = authorPersistence.NewRoleRepositoryImpl()
		orderRepo              = orderPersistence.NewOrderRepositoryImpl()
		orderItemRepo          = orderPersistence.NewOrderItemRepositoryImpl()
		invoiceRepo            = orderPersistence.NewInvoiceRepositoryImpl()
		invoiceItemRepo        = orderPersistence.NewInvoiceItemRepositoryImpl()
		shipmentRepo           = shipmentPersistence.NewShipmentRepositoryImpl()
	)

	var (
		tokenSvc    = authenSvc.NewJWTService(cfg.JWT.Secret)
		pwHasher    = authenSvc.NewBcryptHasher()
		tokenHasher = authenSvc.NewSHATokenHasher()
		authMidd    = authenSvc.NewJWTAuthenticator(
			tokenSvc,
			sessionRepo,
			tokenHasher,
			refreshTokenRepo,
		)

		actorSvc = authorSvc.NewActorService(
			accountRepo,
			membershipRepo,
		)
		userDeletionSvc = authenSvc.NewUserDeletionService(
			accountRepo,
			oauthRepo,
			sessionRepo,
			userRepo,
		)
		customerDeletionSvc = customerSvc.NewCustomerDeletionService(
			customerRepo,
			addressRepo,
			cartRepo,
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

		pricingService = orderSvc.NewPricingService(
			addressRepo,
			shopCourierRepo,
			inventoryRepo,
			paymentMethodRepo,
			productRepo,
			infra.ShippingProvider,
			addressShopRepo,
			shopRepo,
		)
	)

	return &Container{
		Logger:             log,
		AuditLogger:        auditLogger,
		CORSAllowedOrigins: cfg.App.CORSAllowedOrigins,
		Authenticator:      authMidd,
		Authorizer:         authorMdwr,
		DBExecutor:         infra.TransactionExecutor,
		DBTransactor:       infra.TransactionProvider,
		GoogleOAuth:        cfg.GoogleOAuth,

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
				infra.TransactionExecutor,
				infra.StorageProvider,
				productRepo,
				inventoryRepo,
				productImageRepo,
				shopRepo,
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
			userRepo,
			actorSvc,
			oauthRepo,
		),
		Logout: *authenUsecase.NewLogoutUsecase(
			infra.TransactionProvider,
			refreshTokenRepo,
			sessionRepo,
			auditLogger,
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
				customerRepo,
				auditLogger,
			),
		LoginStaff: *authenUsecase.
			NewLoginStaffUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				tokenHasher,
				tokenSvc,
				sessionRepo,
				refreshTokenRepo,
				staffRepo,
				membershipRepo,
				auditLogger,
			),

		FindStaff: *staffUsecase.
			NewFindStaffUsecase(
				infra.TransactionExecutor,
				staffRepo,
			),
		CreateStaff: *staffUsecase.
			NewCreateStaffUsecase(
				staffRepo,
				userRepo,
				infra.TransactionExecutor,
				infra.TransactionProvider,
				auditLogger,
			),
		AddStaffAccount: *staffUsecase.
			NewAddStaffAccountUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				userRepo,
				membershipRepo,
				roleRepo,
				auditLogger,
			),

		RegisterCustomer: *authenUsecase.
			NewRegisterCustomerUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				userRepo,
				customerRepo,
				challengeRepo,
				otpGen,
				mailSender,
				auditLogger,
			),
		VerifyAccount: *authenUsecase.
			NewVerifyAccountUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				pwHasher,
				tokenHasher,
				userRepo,
				customerRepo,
				membershipRepo,
				challengeRepo,
				tokenSvc,
				sessionRepo,
				refreshTokenRepo,
				auditLogger,
			),
		GetAccount: *authenUsecase.
			NewGetAccountUsecase(
				accountRepo,
				infra.TransactionExecutor,
			),
		AuthenticateOAuth: *authenUsecase.
			NewAuthenticateOAuthUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				oauthRepo,
				userRepo,
				customerRepo,
				tokenHasher,
				tokenSvc,
				sessionRepo,
				refreshTokenRepo,
				auditLogger,
			),
		RequestPasswordReset: *authenUsecase.
			NewRequestPasswordResetUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				challengeRepo,
				pwHasher,
				otpGen,
				mailSender,
				auditLogger,
			),
		VerifyPasswordReset: *authenUsecase.
			NewVerifyPasswordResetUsecase(
				infra.TransactionExecutor,
				challengeRepo,
				pwHasher,
				auditLogger,
			),
		ResetPassword: *authenUsecase.
			NewResetPasswordUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				sessionRepo,
				challengeRepo,
				pwHasher,
				auditLogger,
			),
		DeleteCustomerAccount: *customerUsecase.
			NewDeleteCustomerAccountUsecase(
				infra.TransactionProvider,
				userDeletionSvc,
				customerDeletionSvc,
				auditLogger,
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
				pricingService,
			),

		ListLocations: *locationUsecase.
			NewListLocationUsecase(
				infra.ShippingProvider,
				infra.TransactionExecutor,
			),

		GetUser: *userUsecase.
			NewGetUserUsecase(
				userRepo,
				infra.TransactionExecutor,
			),
		GetCurrentProfile: *userUsecase.NewGetCurrentProfileUsecase(
			infra.TransactionExecutor,
			accountRepo,
			customerRepo,
			staffRepo,
			sessionRepo,
		),
		UpdateCurrentProfile: *userUsecase.NewUpdateCurrentProfileUsecase(
			infra.TransactionExecutor,
			infra.TransactionProvider,
			accountRepo,
			customerRepo,
			staffRepo,
			userRepo,
		),

		FindCustomers: *customerUsecase.
			NewFindCustomersUsecase(
				infra.TransactionExecutor,
				customerRepo,
			),

		ListUserAddresses: *addressUsecase.
			NewListCustomerAddressesUsecase(
				addressRepo,
				infra.TransactionExecutor,
			),
		CreateUserAddress: *addressUsecase.
			NewSaveCustomerAddressUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				addressRepo,
			),
		DeleteUserAddress: *addressUsecase.
			NewDeleteCustomerAddressUsecase(
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
		SaveShop: *shopUsecase.
			NewSaveShopUsecase(
				shopRepo,
				slugGen,
				infra.TransactionExecutor,
			),

		GetShopAddresses: *shopUsecase.
			NewGetShopAddressesUsecase(
				addressShopRepo,
				infra.TransactionExecutor,
			),
		GetShopCouriers: *shopUsecase.
			NewGetShopCouriersUsecase(
				shopCourierRepo,
				infra.TransactionExecutor,
			),
		GetShopProducts: *shopUsecase.
			NewGetShopProductsUsecase(
				inventoryRepo,
				productRepo,
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
		ProcessPaymentWebhook: *paymentUsecase.
			NewProcessPaymentWebhookUsecase(
				paymentRepo,
				paymentAccRepo,
				paymentEventRepo,
				orderRepo,
				orderItemRepo,
				inventoryRepo,
				infra.PaymentGateway,
				infra.TransactionProvider,
				infra.TransactionExecutor,
			),
		ProcessManualPayment: *paymentUsecase.
			NewProcessManualPaymentUsecase(
				paymentRepo,
				paymentAccRepo,
				paymentEventRepo,
				orderRepo,
				orderItemRepo,
				inventoryRepo,
				infra.TransactionProvider,
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
				infra.ShippingProvider,
				infra.TransactionExecutor,
			),

		CreateOrder: *orderUsecase.
			NewCreateOrderUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				accountRepo,
				orderRepo,
				orderItemRepo,
				invoiceRepo,
				invoiceItemRepo,
				paymentRepo,
				paymentMethodRepo,
				paymentAccRepo,
				paymentEventRepo,
				paymentInstructionRepo,
				inventoryRepo,
				cartRepo,
				userRepo,
				infra.PaymentGateway,
				pricingService,
			),
		FindOrders: *orderUsecase.
			NewFindOrdersUsecase(
				infra.TransactionExecutor,
				orderRepo,
				orderItemRepo,
				paymentRepo,
				shipmentRepo,
			),
		GetOrder: *orderUsecase.
			NewGetOrderUsecase(
				infra.TransactionExecutor,
				orderRepo,
				orderItemRepo,
				paymentRepo,
				shipmentRepo,
			),
		UpdateOrderStatus: *orderUsecase.
			NewUpdateOrderStatusUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				orderRepo,
				orderItemRepo,
				shipmentRepo,
				addressRepo,
				addressShopRepo,
				infra.LogisticsProvider,
			),

		FindAuditLogs: *auditUsecase.NewFindAuditLogsUsecase(
			infra.TransactionExecutor,
			auditLogRepo,
		),
		GetAuditLog: *auditUsecase.NewGetAuditLogUsecase(
			infra.TransactionExecutor,
			auditLogRepo,
		),
		DeleteAuditLogs: *auditUsecase.NewDeleteAuditLogsUsecase(
			infra.TransactionExecutor,
			auditLogRepo,
		),

		WAFAutoBanEnabled: cfg.WAF.AutoBanEnabled,
		Limiter:           applimiter.NewInMemorySlidingWindowLimiter(10*time.Second, 30),
		ListRules: *secPolicyUsecase.NewListRulesUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		CreateRule: *secPolicyUsecase.NewCreateRuleUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		ToggleRule: *secPolicyUsecase.NewToggleRuleUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		UpdateRule: *secPolicyUsecase.NewUpdateRuleUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		DeleteRule: *secPolicyUsecase.NewDeleteRuleUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		GetIPConfig: *secPolicyUsecase.NewGetIPConfigUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		UpdateIPAction: *secPolicyUsecase.NewUpdateIPActionUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		GetFilters: *secPolicyUsecase.NewGetFiltersUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		UpdateFilter: *secPolicyUsecase.NewUpdateFilterUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
		InspectPayload: *secPolicyUsecase.NewInspectPayloadUsecase(
			infra.TransactionExecutor,
			secPolicyRepo,
		),
	}
}
