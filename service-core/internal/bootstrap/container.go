package bootstrap

import (
	"time"

	applimiter "service-core/internal/common/limiter"
	applogger "service-core/internal/common/logger"

	"service-core/internal/infra/shipping"
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

	threatIntelProvider "service-core/internal/modules/threat_intel/infra/provider"
	threatIntelUsecase "service-core/internal/modules/threat_intel/usecase"

	addressPersistence "service-core/internal/modules/address/infra/persistence"
	analyticsPersistence "service-core/internal/modules/analytics/infra/persistence"
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
	authorSvc "service-core/internal/modules/authorization/infra/service"
	authorRepo "service-core/internal/modules/authorization/repository"
	customerSvc "service-core/internal/modules/customer/infra/service"
	orderSvc "service-core/internal/modules/order/infra/service"

	addressUsecase "service-core/internal/modules/address/usecase"
	analyticsUsecase "service-core/internal/modules/analytics/usecase"
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

	intelligencelayer "service-core/internal/infra/intelligence_layer"
	paymentgateway "service-core/internal/infra/payment-gateway"
	paymentRepo "service-core/internal/modules/payment/repository"
)

type Container struct {
	Logger               applogger.Logger
	AuditLogger          applogger.AuditLogger
	CORSAllowedOrigins   []string
	Authenticator        authenRepo.Authenticator
	Authorizer           authorRepo.Authorizer
	DBExecutor           transaction.Executor
	DBTransactor         transaction.Transactor
	GoogleOAuth          appconfig.GoogleOAuthConfig
	paymentMethodRepo    paymentRepo.PaymentMethodRepository
	paymentGateway       paymentgateway.Provider
	IntelligenceProvider intelligencelayer.Provider

	FindProducts     productUsecase.FindProductsUsecase
	GetProduct       productUsecase.GetProductUsecase
	SaveProduct      productUsecase.SaveProductUsecase
	DeleteProduct    productUsecase.DeleteProductUsecase
	AddProductImages productUsecase.AddProductImagesUsecase
	GetProductStats  productUsecase.GetProductStatsUsecase
	CreateInventory  inventoryUsecase.CreateInventoryUsecase
	UpdateInventory  inventoryUsecase.UpdateInventoryUsecase
	DeleteInventory  inventoryUsecase.DeleteInventoryUsecase

	Me                    authenUsecase.MeUsecase
	LoginCustomer         authenUsecase.LoginCustomerUsecase
	LoginStaff            authenUsecase.LoginStaffUsecase
	RegisterCustomer      authenUsecase.RegisterCustomerUsecase
	VerifyAccount         authenUsecase.VerifyAccountUsecase
	GetAccount            authenUsecase.GetAccountUsecase
	Logout                authenUsecase.LogoutUsecase
	AuthenticateOAuth     authenUsecase.AuthenticateOAuthUsecase
	RequestPasswordReset  authenUsecase.RequestPasswordResetUsecase
	VerifyPasswordReset   authenUsecase.VerifyPasswordResetUsecase
	ResetPassword         authenUsecase.ResetPasswordUsecase
	DeleteCustomerAccount customerUsecase.DeleteCustomerAccountUsecase

	FindStaff          staffUsecase.FindStaffUsecase
	CreateStaff        staffUsecase.CreateStaffUsecase
	AddStaffAccount    staffUsecase.AddStaffAccountUsecase
	ListStaffAccounts  staffUsecase.ListStaffAccountsUsecase
	UpdateStaff        staffUsecase.UpdateStaffUsecase
	DeleteStaff        staffUsecase.DeleteStaffUsecase
	RemoveStaffAccount staffUsecase.RemoveStaffAccountUsecase

	ListStaffPermissions  staffUsecase.ListStaffPermissionsUsecase
	SaveStaffPermission   staffUsecase.SaveStaffPermissionUsecase
	DeleteStaffPermission staffUsecase.DeleteStaffPermissionUsecase

	GetCart          cartUsecase.GetCartUsecase
	AddItem          cartUsecase.AddItemUsecase
	AddCustomItem    cartUsecase.AddCustomItemUsecase
	UpdateItem       cartUsecase.UpdateItemUsecase
	RemoveItem       cartUsecase.RemoveItemUsecase
	RemoveCustomItem cartUsecase.RemoveCustomItemUsecase
	ChangeItemShop   cartUsecase.ChangeItemShopUsecase
	Checkout         cartUsecase.CheckoutUsecase

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
	UpdateShopAddress addressUsecase.UpdateShopAddressUsecase
	DeleteShopAddress addressUsecase.DeleteShopAddressUsecase

	FindShops  shopUsecase.FindShopsUsecase
	GetShop    shopUsecase.GetShopUsecase
	SaveShop   shopUsecase.SaveShopUsecase
	DeleteShop shopUsecase.DeleteShopUsecase

	GetShopAddresses shopUsecase.GetShopAddressesUsecase
	GetShopCouriers  shopUsecase.GetShopCouriersUsecase
	GetShopProducts  shopUsecase.GetShopProductsUsecase

	SavePaymentMethod      paymentUsecase.SavePaymentMethodUsecase
	ListPaymentMethod      paymentUsecase.ListPaymentMethodUsecase
	ProcessPaymentWebhook  paymentUsecase.ProcessPaymentWebhookUsecase
	SavePaymentInstruction paymentUsecase.SavePaymentInstructionUsecase
	GetPaymentDetail       paymentUsecase.GetPaymentDetailUsecase
	CheckPaymentStatus     paymentUsecase.CheckPaymentStatusUsecase
	SyncPendingPayments    paymentUsecase.SyncPendingPaymentsUsecase
	ExpirePastDuePayments  paymentUsecase.ExpirePastDuePaymentsUsecase
	SyncPaymentMethods     paymentUsecase.SyncPaymentMethodsUsecase
	ProcessOrderRefund     paymentUsecase.ProcessOrderRefundUsecase

	ListAllCouriers      courierUsecase.ListCouriersUsecase
	ConfigureShopCourier courierUsecase.ConfigureShopCourierUsecase

	EstimateShippingOptions shipmentUsecase.EstimateShippingOptionsUsecase
	UpdateShipmentStatus    shipmentUsecase.UpdateShipmentStatusUsecase
	UpdateShipment          shipmentUsecase.UpdateShipmentUsecase

	CreateOrder             orderUsecase.CreateOrderUsecase
	FindOrders              orderUsecase.FindOrdersUsecase
	GetOrder                orderUsecase.GetOrderUsecase
	UpdateOrderStatus       orderUsecase.UpdateOrderStatusUsecase
	DispatchShopShipment    orderUsecase.DispatchShopShipmentUsecase
	GetOrderTracking        orderUsecase.GetOrderTrackingUsecase
	ExpireUnfulfilledOrders orderUsecase.ExpireUnfulfilledOrdersUsecase

	FindAuditLogs   auditUsecase.FindAuditLogsUsecase
	GetAuditLog     auditUsecase.GetAuditLogUsecase
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

	AnalyzeIP threatIntelUsecase.AnalyzeIPUsecase
	GetGeoIP  threatIntelUsecase.GetGeoIPUsecase

	GetOrderMetrics     analyticsUsecase.GetOrderMetricsUsecase
	GetPaymentMetrics   analyticsUsecase.GetPaymentMetricsUsecase
	GetShipmentMetrics  analyticsUsecase.GetShipmentMetricsUsecase
	GetInventoryMetrics analyticsUsecase.GetInventoryMetricsUsecase
	GetProductMetrics   analyticsUsecase.GetProductMetricsUsecase
	GetDemandForecast   analyticsUsecase.GetDemandForecastUsecase
	GetStockoutRisks    analyticsUsecase.GetStockoutRisksUsecase
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
		productRepo             = productPersistence.NewProductRepository()
		productImageRepo        = productPersistence.NewProductImageRepository()
		productPerformanceRepo  = productPersistence.NewProductPerformanceRepository()
		productStockHistoryRepo = productPersistence.NewProductStockHistoryRepository()
		inventoryRepo           = inventoryPersistence.NewInventoryRepository()
		secPolicyRepo           = secPolicyPersistence.NewSecurityPolicyRepository()
		accountRepo             = authenPersistence.NewAccountRepository()
		challengeRepo           = authenPersistence.NewChallengeRepository()
		oauthRepo               = authenPersistence.NewOAuthConnectionRepository()
		sessionRepo             = authenPersistence.NewSessionRepositoryImpl()
		refreshTokenRepo        = authenPersistence.NewRefreshTokenRepositoryImpl()
		cartRepo                = cartPersistence.NewCartRepositoryImpl()
		userRepo                = userPersistence.NewUserRepositoryImpl()
		addressRepo             = addressPersistence.NewCustomerAddressRepositoryImpl()
		addressShopRepo         = addressPersistence.NewShopAddressRepositoryImpl()
		paymentRepo             = paymentPersistence.NewPaymentRepositoryImpl()
		paymentMethodRepo       = paymentPersistence.NewPaymentMethodRepository()
		paymentEventRepo        = paymentPersistence.NewPaymentEventRepositoryImpl()
		paymentInstructionRepo  = paymentPersistence.NewPaymentInstructionRepositoryImpl()
		paymentChannelDataRepo  = paymentPersistence.NewPaymentChannelDataRepositoryImpl()
		paymentWebhookEventRepo = paymentPersistence.NewPaymentWebhookEventRepositoryImpl()
		shopRepo                = shopPersistence.NewShopRepositoryImpl()
		courierRepo             = courierPersistence.NewCourierRepositoryImpl()
		shopCourierRepo         = courierPersistence.NewShopCourierRepositoryImpl()
		staffRepo               = staffPersistence.NewStaffRepositoryImpl()
		customerRepo            = customerPersistence.NewCustomerRepositoryImpl()
		membershipRepo          = authorPersistence.NewStaffMembershipRepositoryImpl()
		staffPermRepo           = authorPersistence.NewStaffPermissionRepositoryImpl()
		roleRepo                = authorPersistence.NewRoleRepositoryImpl()
		orderRepo               = orderPersistence.NewOrderRepositoryImpl()
		orderItemRepo           = orderPersistence.NewOrderItemRepositoryImpl()
		invoiceRepo             = orderPersistence.NewInvoiceRepositoryImpl()
		invoiceItemRepo         = orderPersistence.NewInvoiceItemRepositoryImpl()
		shipmentRepo            = shipmentPersistence.NewShipmentRepositoryImpl()
		shipmentEventRepo       = shipmentPersistence.NewShipmentEventRepositoryImpl()
		threatIntelRepo         = threatIntelProvider.NewThreatIntelProvider(cfg.WAF)
		analyticsRepo           = analyticsPersistence.NewAnalyticsRepositoryImpl()
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
			staffPermRepo,
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
			cartRepo,
			shopCourierRepo,
			inventoryRepo,
			paymentMethodRepo,
			productRepo,
			infra.ShippingProvider,
			addressShopRepo,
			shopRepo,
		)
	)

	var intelligenceProvider intelligencelayer.Provider
	if cfg.IntelligenceLayer.Enabled {
		intelligenceProvider = intelligencelayer.NewClient(
			cfg.IntelligenceLayer.BaseURL,
			time.Duration(cfg.IntelligenceLayer.TimeoutMS)*time.Millisecond,
			log,
		)
	}

	processPaymentWebhook := *paymentUsecase.NewProcessPaymentWebhookUsecase(
		paymentRepo,
		paymentEventRepo,
		paymentWebhookEventRepo,
		orderRepo,
		orderItemRepo,
		inventoryRepo,
		infra.PaymentGateway,
		auditLogger,
		infra.TransactionProvider,
		infra.TransactionExecutor,
		intelligenceProvider,
	)

	c := &Container{
		Logger:               log,
		AuditLogger:          auditLogger,
		CORSAllowedOrigins:   cfg.App.CORSAllowedOrigins,
		Authenticator:        authMidd,
		Authorizer:           authorMdwr,
		DBExecutor:           infra.TransactionExecutor,
		DBTransactor:         infra.TransactionProvider,
		GoogleOAuth:          cfg.GoogleOAuth,
		paymentMethodRepo:    paymentMethodRepo,
		paymentGateway:       infra.PaymentGateway,
		IntelligenceProvider: intelligenceProvider,

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
				productPerformanceRepo,
			),
		SaveProduct: *productUsecase.
			NewSaveProductUsecase(
				infra.TransactionProvider,
				productRepo,
				slugGen,
				productPerformanceRepo,
			),
		GetProductStats: *productUsecase.
			NewGetProductStatsUsecase(
				productPerformanceRepo,
				productImageRepo,
				infra.StorageProvider,
				infra.TransactionExecutor,
			),
		DeleteProduct: *productUsecase.
			NewDeleteProductUsecase(
				productRepo,
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
				productStockHistoryRepo,
			),
		UpdateInventory: *inventoryUsecase.
			NewUpdateInventoryUsecase(inventoryRepo,
				infra.TransactionExecutor,
				productStockHistoryRepo,
			),
		DeleteInventory: *inventoryUsecase.
			NewDeleteInventoryUsecase(inventoryRepo,
				infra.TransactionExecutor,
				productStockHistoryRepo,
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
				staffRepo,
				membershipRepo,
				roleRepo,
				auditLogger,
			),
		ListStaffAccounts: *staffUsecase.
			NewListStaffAccountsUsecase(
				infra.TransactionExecutor,
				staffRepo,
				membershipRepo,
				auditLogger,
			),
		UpdateStaff: *staffUsecase.
			NewUpdateStaffUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				staffRepo,
				membershipRepo,
				auditLogger,
			),
		DeleteStaff: *staffUsecase.
			NewDeleteStaffUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				staffRepo,
				membershipRepo,
				userDeletionSvc,
				auditLogger,
			),
		RemoveStaffAccount: *staffUsecase.
			NewRemoveStaffAccountUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				staffRepo,
				membershipRepo,
				accountRepo,
				sessionRepo,
				auditLogger,
			),
		ListStaffPermissions: *staffUsecase.NewListStaffPermissionsUsecase(
			infra.TransactionExecutor,
			staffPermRepo,
		),
		SaveStaffPermission: *staffUsecase.NewSaveStaffPermissionUsecase(
			infra.TransactionProvider,
			staffPermRepo,
		),
		DeleteStaffPermission: *staffUsecase.NewDeleteStaffPermissionUsecase(
			infra.TransactionProvider,
			staffPermRepo,
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
				shopRepo,
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
				shopRepo,
			),
		AddCustomItem: *cartUsecase.
			NewAddCustomItemUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				cartRepo,
				shopRepo,
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
		RemoveCustomItem: *cartUsecase.
			NewRemoveCustomItemUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				cartRepo,
			),
		ChangeItemShop: *cartUsecase.
			NewChangeItemShopUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				cartRepo,
				shopRepo,
				inventoryRepo,
			),
		Checkout: *cartUsecase.
			NewCheckoutUsecase(
				infra.TransactionExecutor,
				pricingService,
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
		UpdateShopAddress: *addressUsecase.
			NewUpdateShopAddressUsecase(
				addressShopRepo,
				infra.TransactionExecutor,
				infra.TransactionProvider,
			),
		DeleteShopAddress: *addressUsecase.
			NewDeleteShopAddressUsecase(
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
		DeleteShop: *shopUsecase.
			NewDeleteShopUsecase(
				shopRepo,
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

		SavePaymentMethod: *paymentUsecase.
			NewSavePaymentMethodUsecase(
				paymentMethodRepo,
				infra.TransactionExecutor,
			),
		ListPaymentMethod: *paymentUsecase.
			NewListPaymentMethodUsecase(
				paymentMethodRepo,
				infra.TransactionExecutor,
			),
		ProcessPaymentWebhook: processPaymentWebhook,
		SavePaymentInstruction: *paymentUsecase.
			NewSavePaymentInstructionUsecase(
				paymentMethodRepo,
				paymentInstructionRepo,
				infra.TransactionExecutor,
			),
		GetPaymentDetail: *paymentUsecase.
			NewGetPaymentDetailUsecase(
				infra.TransactionExecutor,
				orderRepo,
				invoiceRepo,
				paymentRepo,
				paymentMethodRepo,
				paymentInstructionRepo,
				paymentChannelDataRepo,
			),
		CheckPaymentStatus: *paymentUsecase.
			NewCheckPaymentStatusUsecase(
				orderRepo,
				paymentRepo,
				infra.PaymentGateway,
				&processPaymentWebhook,
				infra.TransactionExecutor,
			),
		SyncPendingPayments: *paymentUsecase.
			NewSyncPendingPaymentsUsecase(
				paymentRepo,
				infra.PaymentGateway,
				&processPaymentWebhook,
				infra.TransactionExecutor,
				log,
				time.Duration(cfg.PaymentSync.LookbackHours)*time.Hour,
				infra.TransactionProvider,
				orderRepo,
				orderItemRepo,
				inventoryRepo,
			),
		ExpirePastDuePayments: *paymentUsecase.
			NewExpirePastDuePaymentsUsecase(
				paymentRepo,
				infra.PaymentGateway,
				infra.TransactionExecutor,
				infra.TransactionProvider,
				orderRepo,
				orderItemRepo,
				inventoryRepo,
				log,
				cfg.PaymentExpiry.BatchSize,
				cfg.PaymentExpiry.Concurrency,
			),
		SyncPaymentMethods: *paymentUsecase.
			NewSyncPaymentMethodsUsecase(
				paymentMethodRepo,
				infra.TransactionExecutor,
				infra.PaymentGateway,
			),
		ProcessOrderRefund: *paymentUsecase.NewProcessOrderRefundUsecase(
			paymentRepo,
			infra.PaymentGateway,
			infra.TransactionExecutor,
			infra.TransactionProvider,
			log,
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
				intelligenceProvider,
			),
		UpdateShipmentStatus: *shipmentUsecase.
			NewUpdateShipmentStatusUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				shipmentRepo,
				shipmentEventRepo,
				orderRepo,
			),
		UpdateShipment: *shipmentUsecase.
			NewUpdateShipmentUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				shipmentRepo,
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
				paymentEventRepo,
				paymentInstructionRepo,
				paymentChannelDataRepo,
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
				paymentChannelDataRepo,
				shipmentRepo,
				addressRepo,
			),
		GetOrder: *orderUsecase.
			NewGetOrderUsecase(
				infra.TransactionExecutor,
				orderRepo,
				orderItemRepo,
				paymentRepo,
				paymentChannelDataRepo,
				shipmentRepo,
				shipmentEventRepo,
			),
		UpdateOrderStatus: *orderUsecase.
			NewUpdateOrderStatusUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				orderRepo,
				orderItemRepo,
				inventoryRepo,
				paymentRepo,
				productRepo,
				shipmentRepo,
				addressRepo,
				addressShopRepo,
				infra.LogisticsProvider,
				auditLogger,
			),
		DispatchShopShipment: *orderUsecase.
			NewDispatchShopShipmentUsecase(
				infra.TransactionExecutor,
				infra.TransactionProvider,
				orderRepo,
				orderItemRepo,
				productRepo,
				shipmentRepo,
				addressRepo,
				addressShopRepo,
				infra.LogisticsProvider,
				auditLogger,
			),
		GetOrderTracking: *orderUsecase.
			NewGetOrderTrackingUsecase(
				infra.TransactionExecutor,
				orderRepo,
				shipmentRepo,
				shipmentEventRepo,
				infra.LogisticsProvider,
				addressRepo,
				shipping.NewTrackingCache(shipping.DefaultTrackingCacheTTL),
			),
		ExpireUnfulfilledOrders: *orderUsecase.NewExpireUnfulfilledOrdersUsecase(
			orderRepo,
			orderItemRepo,
			inventoryRepo,
			paymentUsecase.NewProcessOrderRefundUsecase(
				paymentRepo,
				infra.PaymentGateway,
				infra.TransactionExecutor,
				infra.TransactionProvider,
				log,
			),
			infra.TransactionExecutor,
			infra.TransactionProvider,
			log,
			auditLogger,
			100,
			5,
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
		AnalyzeIP: *threatIntelUsecase.NewAnalyzeIPUsecase(threatIntelRepo),
		GetGeoIP:  *threatIntelUsecase.NewGetGeoIPUsecase(threatIntelRepo),

		GetOrderMetrics:     *analyticsUsecase.NewGetOrderMetricsUsecase(infra.TransactionExecutor, analyticsRepo),
		GetPaymentMetrics:   *analyticsUsecase.NewGetPaymentMetricsUsecase(infra.TransactionExecutor, analyticsRepo),
		GetShipmentMetrics:  *analyticsUsecase.NewGetShipmentMetricsUsecase(infra.TransactionExecutor, analyticsRepo),
		GetInventoryMetrics: *analyticsUsecase.NewGetInventoryMetricsUsecase(infra.TransactionExecutor, analyticsRepo),
		GetProductMetrics:   *analyticsUsecase.NewGetProductMetricsUsecase(infra.TransactionExecutor, analyticsRepo),
		GetDemandForecast:   *analyticsUsecase.NewGetDemandForecastUsecase(infra.TransactionExecutor, analyticsRepo, intelligenceProvider),
		GetStockoutRisks:    *analyticsUsecase.NewGetStockoutRisksUsecase(infra.TransactionExecutor, analyticsRepo, intelligenceProvider),
	}

	return c
}
