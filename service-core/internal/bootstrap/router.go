package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	appmiddleware "service-core/internal/common/middleware"

	authendomain "service-core/internal/modules/authentication/domain"
	authorzDomain "service-core/internal/modules/authorization/domain"

	addressH "service-core/internal/modules/address/delivery/http"
	analyticsH "service-core/internal/modules/analytics/delivery/http"
	auditH "service-core/internal/modules/audit/delivery/http"
	authH "service-core/internal/modules/authentication/delivery/http"
	cartH "service-core/internal/modules/cart/delivery/http"
	courierH "service-core/internal/modules/courier/delivery/http"
	customerH "service-core/internal/modules/customer/delivery/http"
	inventoryH "service-core/internal/modules/inventory/delivery/http"
	locationH "service-core/internal/modules/location/delivery/http"
	orderH "service-core/internal/modules/order/delivery/http"
	paymentH "service-core/internal/modules/payment/delivery/http"
	productH "service-core/internal/modules/product/delivery/http"
	secPolicyH "service-core/internal/modules/security_policy/delivery/http"
	shipmentH "service-core/internal/modules/shipment/delivery/http"
	shopH "service-core/internal/modules/shop/delivery/http"
	staffH "service-core/internal/modules/staff/delivery/http"
	threatIntelH "service-core/internal/modules/threat_intel/delivery/http"
	userH "service-core/internal/modules/user/delivery/http"

	"github.com/go-chi/chi/v5"
)

// RouteChains encapsulates all pre-built middleware chains for
// different routing policies.
type RouteChains struct {
	Core           func(apphttp.AppHandler) http.HandlerFunc
	CoreAuth       func(apphttp.AppHandler) http.HandlerFunc
	StaffOnly      func(apphttp.AppHandler) http.HandlerFunc
	StaffAdminOnly func(apphttp.AppHandler) http.HandlerFunc
	CustomerOnly   func(apphttp.AppHandler) http.HandlerFunc
}

// NewRouteChains builds and returns the route chains using
// the provided Container.
func NewRouteChains(c *Container) *RouteChains {
	buildChain := func(
		extra ...appmiddleware.Middleware,
	) func(apphttp.AppHandler) http.HandlerFunc {
		base := []appmiddleware.Middleware{
			appmiddleware.CORS(c.CORSAllowedOrigins),
			appmiddleware.Recovery(c.Logger),
			appmiddleware.Logging(c.Logger),
			appmiddleware.Response(),
			appmiddleware.WAF(
				&c.InspectPayload,
				&c.UpdateIPAction,
				c.AuditLogger,
				c.Limiter,
				c.Logger,
				c.WAFAutoBanEnabled,
			),
		}

		mws := append(base, extra...)

		return func(h apphttp.AppHandler) http.HandlerFunc {
			return appmiddleware.Chain(h, mws...)
		}
	}

	return &RouteChains{
		Core: buildChain(),
		CoreAuth: buildChain(
			c.Authenticator.RequireMultiAuth(
				c.DBExecutor,
				c.DBTransactor,
				appcookie.CookieCustomer,
				appcookie.CookieStaff,
			),
			c.Authorizer.RequireAccountType(
				authendomain.AccountTypeStaff,
				authendomain.AccountTypeCustomer,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
		),
		StaffOnly: buildChain(
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				c.DBTransactor,
				appcookie.CookieStaff,
			),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeStaff),
			c.Authorizer.LoadActor(c.DBExecutor),
		),
		StaffAdminOnly: buildChain(
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				c.DBTransactor,
				appcookie.CookieStaff,
			),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeStaff),
			c.Authorizer.RequireStaffRole(authorzDomain.RoleStaffAdmin),
			c.Authorizer.LoadActor(c.DBExecutor),
		),
		CustomerOnly: buildChain(
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				c.DBTransactor,
				appcookie.CookieCustomer,
			),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeCustomer),
			c.Authorizer.LoadActor(c.DBExecutor),
		),
	}
}

func NewRouter(c *Container) *chi.Mux {
	chains := NewRouteChains(c)

	var (
		productHandler = productH.NewProductHandler(
			&c.FindProducts,
			&c.GetProduct,
			&c.SaveProduct,
			&c.DeleteProduct,
			&c.AddProductImages,
			&c.GetProductStats,
		)

		inventoryHandler = inventoryH.NewInventoryHandler(
			&c.CreateInventory,
			&c.UpdateInventory,
			&c.DeleteInventory,
		)

		authHandler = authH.NewAuthHandler(
			&c.Me,
			&c.Logout,
			&c.LoginCustomer,
			&c.LoginStaff,
			&c.RegisterCustomer,
			&c.VerifyAccount,
			&c.GetAccount,
			&c.AuthenticateOAuth,
			&c.RequestPasswordReset,
			&c.VerifyPasswordReset,
			&c.ResetPassword,
			&c.DeleteCustomerAccount,
			c.GoogleOAuth,
		)

		staffHandler = staffH.NewStaffHandler(
			&c.AddStaffAccount,
			&c.CreateStaff,
			&c.FindStaff,
		)

		customerHandler = customerH.NewCustomerHandler(
			&c.FindCustomers,
		)

		cartHandler = cartH.NewCartHandler(
			&c.AddItem,
			&c.AddCustomItem,
			&c.GetCart,
			&c.UpdateItem,
			&c.RemoveItem,
			&c.RemoveCustomItem,
			&c.ChangeItemShop,
			&c.Checkout,
		)

		locationHandler = locationH.NewLocationHandler(
			&c.ListLocations,
		)

		userHandler = userH.NewUserHandler(
			&c.GetUser,
			&c.GetCurrentProfile,
			&c.UpdateCurrentProfile,
		)

		addressHandler = addressH.NewAddressHandler(
			&c.ListUserAddresses,
			&c.CreateUserAddress,
			&c.DeleteUserAddress,
			&c.ListShopAddresses,
			&c.SaveShopAddress,
			&c.UpdateShopAddress,
			&c.DeleteShopAddress,
		)

		paymentHandler = paymentH.NewPaymentHandler(
			&c.SavePaymentMethod,
			&c.ListPaymentMethod,
			&c.ProcessPaymentWebhook,
			&c.SavePaymentInstruction,
			&c.GetPaymentDetail,
			&c.CheckPaymentStatus,
		)

		shopHandler = shopH.NewShopHandler(
			&c.FindShops,
			&c.GetShop,
			&c.SaveShop,
			&c.DeleteShop,
			&c.GetShopAddresses,
			&c.GetShopCouriers,
			&c.GetShopProducts,
		)

		courierHandler = courierH.NewCourierHandler(
			&c.ListAllCouriers,
			&c.ConfigureShopCourier,
		)

		shipmentHandler = shipmentH.NewShipmentHandler(
			&c.EstimateShippingOptions,
			&c.UpdateShipmentStatus,
			&c.UpdateShipment,
		)

		orderHandler = orderH.NewOrderHandler(
			&c.FindOrders,
			&c.GetOrder,
			&c.CreateOrder,
			&c.UpdateOrderStatus,
			&c.GetOrderTracking,
		)

		auditHandler = auditH.NewAuditHandler(
			&c.FindAuditLogs,
			&c.GetAuditLog,
			&c.DeleteAuditLogs,
		)

		secPolicyHandler = secPolicyH.NewSecurityPolicyHandler(
			&c.ListRules,
			&c.CreateRule,
			&c.ToggleRule,
			&c.UpdateRule,
			&c.DeleteRule,
			&c.GetIPConfig,
			&c.UpdateIPAction,
			&c.GetFilters,
			&c.UpdateFilter,
		)

		threatIntelHandler = threatIntelH.NewThreatIntelHandler(
			&c.AnalyzeIP,
			&c.GetGeoIP,
		)

		analyticsHandler = analyticsH.NewAnalyticsHandler(
			&c.GetOrderMetrics,
			&c.GetPaymentMetrics,
			&c.GetShipmentMetrics,
			&c.GetInventoryMetrics,
			&c.GetProductMetrics,
		)
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"healthy"}`))
		})

		r.Route("/products", func(r chi.Router) {
			r.Get("/", chains.Core(productHandler.FindProducts))
			r.Post("/", chains.StaffAdminOnly(productHandler.SaveProduct))
			r.Get("/stats", chains.StaffAdminOnly(productHandler.GetProductStats))

			r.Get("/{slug}", chains.Core(productHandler.GetProduct))

			r.Route("/id/{id}", func(r chi.Router) {
				r.Delete("/", chains.StaffAdminOnly(productHandler.DeleteProduct))
				r.Post("/images", chains.StaffAdminOnly(productHandler.AddProductImages))
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", chains.Core(authHandler.SignInEmail))
			r.Post("/staff/signin", chains.Core(authHandler.SignInStaffEmail))
			r.Post("/signup", chains.Core(authHandler.SignUpAccount))
			r.Post("/verify", chains.Core(authHandler.VerifyAccount))

			r.Post("/forgot-password", chains.Core(authHandler.ForgotPasswordCustomer))
			r.Post("/staff/forgot-password", chains.Core(authHandler.ForgotPasswordStaff))
			r.Post("/forgot-password/verify", chains.Core(authHandler.VerifyPasswordReset))
			r.Post("/forgot-password/reset", chains.Core(authHandler.ResetPassword))

			r.Post("/logout", chains.CustomerOnly(authHandler.Logout))
			r.Post("/staff/logout", chains.StaffOnly(authHandler.LogoutStaff))
			r.Get("/me", chains.CustomerOnly(authHandler.Me))
			r.Get("/staff/me", chains.StaffOnly(authHandler.Me))

			r.Get("/google/login", chains.Core(authHandler.GoogleLogin))
			r.Get("/google/callback", chains.Core(authHandler.GoogleCallback))
		})

		r.Route("/staff", func(r chi.Router) {
			r.Get("/", chains.StaffAdminOnly(staffHandler.FindStaff))
			r.Post("/", chains.StaffAdminOnly(staffHandler.CreateStaff))
			r.Route("/{staffID}", func(r chi.Router) {
				r.Post("/accounts", chains.StaffAdminOnly(staffHandler.AddStaffAccount))
			})
		})

		r.Route("/customers", func(r chi.Router) {
			r.Get("/", chains.StaffAdminOnly(customerHandler.FindCustomers))
		})

		r.Route("/profile", func(r chi.Router) {
			r.Get("/", chains.CoreAuth(userHandler.GetCurrentProfile))
			r.Put("/", chains.CoreAuth(userHandler.UpdateCurrentProfile))
			r.Get("/staff", chains.StaffOnly(userHandler.GetCurrentProfile))
			r.Put("/staff", chains.StaffOnly(userHandler.UpdateCurrentProfile))
			r.Delete("/", chains.CustomerOnly(authHandler.DeleteAccount))
		})

		r.Route("/carts", func(r chi.Router) {
			r.Get("/", chains.CustomerOnly(cartHandler.GetCart))

			r.Route("/checkout", func(r chi.Router) {
				r.Post("/", chains.CustomerOnly(cartHandler.Checkout))
				r.Post("/calculate", chains.CustomerOnly(cartHandler.CheckoutEstimate))
			})

			r.Route("/items", func(r chi.Router) {
				r.Post("/", chains.CustomerOnly(cartHandler.AddItem))
				r.Put("/{shopID}/{productID}", chains.CustomerOnly(cartHandler.UpdateItem))
				r.Patch("/{cartItemID}/shop", chains.CustomerOnly(cartHandler.ChangeItemShop))
				r.Delete("/{shopID}/{productID}", chains.CustomerOnly(cartHandler.RemoveItem))
				r.Delete("/custom/{cartItemID}", chains.CustomerOnly(cartHandler.RemoveCustomItem))
			})
		})

		r.Route("/couriers", func(r chi.Router) {
			r.Get("/", chains.CoreAuth(courierHandler.ListAllCouriers))
		})

		r.Route("/provinces", func(r chi.Router) {
			r.Get("/", chains.Core(locationHandler.Province))
			r.Get("/{id}/cities", chains.Core(locationHandler.City))
		})
		r.Route("/cities", func(r chi.Router) {
			r.Get("/{id}/districts", chains.Core(locationHandler.District))
		})
		r.Route("/districts", func(r chi.Router) {
			r.Get("/{id}/villages", chains.Core(locationHandler.Village))
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", chains.StaffAdminOnly(userHandler.GetUserByID))
		})

		r.Route("/users/me", func(r chi.Router) {
			r.Get("/", chains.CustomerOnly(userHandler.GetCurrentUser))

			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", chains.CustomerOnly(addressHandler.ListUserAddresses))
				r.Post("/", chains.CustomerOnly(addressHandler.SaveUserAddress))
				r.Delete("/{addressID}", chains.CustomerOnly(addressHandler.DeleteUserAddress))
			})

			r.Route("/orders", func(r chi.Router) {
				r.Get("/", chains.CustomerOnly(orderHandler.ListMyOrders))
				r.Get("/{orderID}", chains.CustomerOnly(orderHandler.GetMyOrder))
				r.Get("/{orderID}/tracking", chains.CustomerOnly(orderHandler.GetMyOrderTracking))
				r.Get("/{orderID}/payment", chains.CustomerOnly(paymentHandler.GetMyOrderPayment))
				r.Post("/{orderID}/payment/check", chains.CustomerOnly(paymentHandler.CheckMyOrderPaymentStatus))
			})
		})

		r.Route("/shops", func(r chi.Router) {
			r.Get("/", chains.Core(shopHandler.FindShops))
			r.Post("/", chains.StaffOnly(shopHandler.SaveShop))

			r.Route("/{shopID}", func(r chi.Router) {
				r.Get("/", chains.Core(shopHandler.GetShopByID))
				r.Delete("/", chains.StaffAdminOnly(shopHandler.DeleteShop))

				r.Route("/addresses", func(r chi.Router) {
					r.Get("/", chains.Core(shopHandler.GetShopAddresses))
					r.Post("/", chains.StaffOnly(addressHandler.CreateShopAddress))
					r.Put("/{addressID}", chains.StaffOnly(addressHandler.UpdateShopAddress))
					r.Delete("/{addressID}", chains.StaffOnly(addressHandler.DeleteShopAddress))
				})

				r.Route("/couriers", func(r chi.Router) {
					r.Get("/", chains.Core(shopHandler.GetShopCouriers))
					r.Post("/", chains.StaffOnly(courierHandler.ConfigureCourierShop))
				})

				r.Route("/products", func(r chi.Router) {
					r.Get("/", chains.Core(shopHandler.GetShopProducts))
					r.Post("/{productID}/inventories", chains.StaffOnly(inventoryHandler.AddInventory))
					r.Put("/{productID}/inventories", chains.StaffOnly(inventoryHandler.UpdateInventory))
					r.Delete("/{productID}/inventories", chains.StaffOnly(inventoryHandler.RemoveInventory))
				})
			})
		})

		r.Route("/midtrans", func(r chi.Router) {
			r.Post("/webhook", chains.Core(paymentHandler.HandleMidtransWebhook))
		})

		r.Route("/payments", func(r chi.Router) {
			r.Route("/methods", func(r chi.Router) {
				r.Get("/", chains.StaffOnly(paymentHandler.ListPaymentMethod))
				r.Patch("/{methodID}", chains.StaffAdminOnly(paymentHandler.UpdatePaymentMethodActive))
				r.Post("/{methodID}/instruction", chains.StaffAdminOnly(paymentHandler.SavePaymentInstruction))
			})
		})

		r.Route("/shipping", func(r chi.Router) {
			r.Post("/cost", chains.CoreAuth(shipmentHandler.EstimateShippingOptions))
		})

		r.Route("/order", func(r chi.Router) {
			r.Post("/", chains.CustomerOnly(orderHandler.CreateOrder))
		})

		r.Route("/orders", func(r chi.Router) {
			r.Get("/", chains.StaffOnly(orderHandler.FindOrders))
			r.Get("/{orderID}", chains.StaffOnly(orderHandler.GetOrder))
			r.Patch("/{orderID}/status", chains.StaffOnly(orderHandler.UpdateOrderStatus))
		})

		r.Route("/shipments", func(r chi.Router) {
			r.Patch("/{shipmentID}/status", chains.StaffOnly(shipmentHandler.UpdateShipmentStatus))
			r.Patch("/{shipmentID}", chains.StaffOnly(shipmentHandler.UpdateShipment))
		})

		r.Route("/api/stats", func(r chi.Router) {
			r.Get("/", chains.StaffAdminOnly(auditHandler.ListAuditLogs))
			r.Delete("/", chains.StaffAdminOnly(auditHandler.DeleteAuditLogs))
			r.Get("/{id}", chains.StaffAdminOnly(auditHandler.GetAuditLog))
			r.Delete("/{id}", chains.StaffAdminOnly(auditHandler.DeleteAuditLogs))
		})

		r.Route("/api/rules", func(r chi.Router) {
			r.Get("/", chains.StaffAdminOnly(secPolicyHandler.ListRules))
			r.Post("/", chains.StaffAdminOnly(secPolicyHandler.CreateRule))
			r.Put("/{id}", chains.StaffAdminOnly(secPolicyHandler.UpdateRule))
			r.Delete("/{id}", chains.StaffAdminOnly(secPolicyHandler.DeleteRule))
		})

		r.Route("/api/ip", func(r chi.Router) {
			r.Get("/", chains.StaffAdminOnly(secPolicyHandler.ListIPConfig))
			r.Post("/", chains.StaffAdminOnly(secPolicyHandler.UpdateIPAction))
		})

		r.Route("/api/filters", func(r chi.Router) {
			r.Get("/", chains.StaffAdminOnly(secPolicyHandler.GetFilters))
			r.Post("/", chains.StaffAdminOnly(secPolicyHandler.UpdateFilter))
		})

		r.Route("/api/analyze", func(r chi.Router) {
			r.Get("/{ip}", chains.StaffAdminOnly(threatIntelHandler.AnalyzeIP))
		})

		r.Route("/api/geo", func(r chi.Router) {
			r.Get("/{ip}", chains.StaffAdminOnly(threatIntelHandler.GetGeolocation))
		})

		r.Route("/analytics", func(r chi.Router) {
			r.Get("/orders", chains.StaffAdminOnly(analyticsHandler.GetOrderMetrics))
			r.Get("/payments", chains.StaffAdminOnly(analyticsHandler.GetPaymentMetrics))
			r.Get("/shipments", chains.StaffAdminOnly(analyticsHandler.GetShipmentMetrics))
			r.Get("/inventory", chains.StaffAdminOnly(analyticsHandler.GetInventoryMetrics))
			r.Get("/products", chains.StaffAdminOnly(analyticsHandler.GetProductMetrics))
		})
	})

	return r
}
