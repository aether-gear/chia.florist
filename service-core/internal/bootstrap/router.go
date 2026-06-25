package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	appmiddleware "service-core/internal/common/middleware"

	authendomain "service-core/internal/modules/authentication/domain"
	authorzDomain "service-core/internal/modules/authorization/domain"

	addressH "service-core/internal/modules/address/delivery/http"
	authH "service-core/internal/modules/authentication/delivery/http"
	cartH "service-core/internal/modules/cart/delivery/http"
	courierH "service-core/internal/modules/courier/delivery/http"
	customerH "service-core/internal/modules/customer/delivery/http"
	inventoryH "service-core/internal/modules/inventory/delivery/http"
	locationH "service-core/internal/modules/location/delivery/http"
	merchantH "service-core/internal/modules/merchant/delivery/http"
	orderH "service-core/internal/modules/order/delivery/http"
	paymentH "service-core/internal/modules/payment/delivery/http"
	productH "service-core/internal/modules/product/delivery/http"
	shipmentH "service-core/internal/modules/shipment/delivery/http"
	shopH "service-core/internal/modules/shop/delivery/http"
	userH "service-core/internal/modules/user/delivery/http"

	"github.com/go-chi/chi/v5"
)

// RouteChains encapsulates all pre-built middleware chains for
// different routing policies.
type RouteChains struct {
	Core              func(apphttp.AppHandler) http.HandlerFunc
	CoreAuth          func(apphttp.AppHandler) http.HandlerFunc
	MerchantOnly      func(apphttp.AppHandler) http.HandlerFunc
	MerchantAdminOnly func(apphttp.AppHandler) http.HandlerFunc
	CustomerOnly      func(apphttp.AppHandler) http.HandlerFunc
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
		}

		mws := append(base, extra...)

		return func(h apphttp.AppHandler) http.HandlerFunc {
			return appmiddleware.Chain(h, mws...)
		}
	}

	return &RouteChains{
		Core: buildChain(),
		CoreAuth: buildChain(
			c.Authenticator.RequireAnyAuth(
				c.DBExecutor,
				appcookie.AccessTokenCookieName,
				appcookie.AccessTokenMerchantCookieName,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
			c.Authorizer.RequireAccountType(
				authendomain.AccountTypeMerchant,
				authendomain.AccountTypeCustomer,
			),
		),
		MerchantOnly: buildChain(
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				appcookie.AccessTokenMerchantCookieName,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeMerchant),
		),
		MerchantAdminOnly: buildChain(
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				appcookie.AccessTokenMerchantCookieName,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeMerchant),
			c.Authorizer.RequireMerchantRole(authorzDomain.RoleMerchantAdmin),
		),
		CustomerOnly: buildChain(
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				appcookie.AccessTokenCookieName,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeCustomer),
		),
	}
}

func NewRouter(c *Container) *chi.Mux {
	chains := NewRouteChains(c)

	var (
		productHandler = productH.NewProductHandler(
			&c.FindProducts,
			&c.GetProduct,
			&c.CreateProduct,
			&c.AddProductImages,
		)

		inventoryHandler = inventoryH.NewInventoryHandler(
			&c.CreateInventory,
		)

		authHandler = authH.NewAuthHandler(
			&c.Me,
			&c.Logout,
			&c.LoginCustomer,
			&c.LoginMerchant,
			&c.RegisterCustomer,
			&c.VerifyAccount,
			&c.GetAccount,
		)

		merchantHandler = merchantH.NewMerchantHandler(
			&c.AddMerchantAccount,
			&c.CreateMerchant,
			&c.FindMerchants,
		)

		customerHandler = customerH.NewCustomerHandler(
			&c.FindCustomers,
		)

		cartHandler = cartH.NewCartHandler(
			&c.AddItem,
			&c.GetCart,
			&c.UpdateItem,
			&c.RemoveItem,
			&c.Checkout,
		)

		locationHandler = locationH.NewLocationHandler(
			&c.ListLocations,
		)

		userHandler = userH.NewUserHandler(
			&c.GetUser,
		)

		addressHandler = addressH.NewAddressHandler(
			&c.ListUserAddresses,
			&c.CreateUserAddress,
			&c.DeleteUserAddress,
			&c.ListShopAddresses,
			&c.SaveShopAddress,
		)

		paymentHandler = paymentH.NewPaymentHandler(
			&c.CreatePaymentAccount,
			&c.ListPaymentAccount,
			&c.CreatePaymentMethod,
			&c.ListPaymentMethod,
			&c.ProcessPaymentWebhook,
			&c.ProcessManualPayment,
		)

		shopHandler = shopH.NewShopHandler(
			&c.FindShops,
			&c.GetShop,
			&c.SaveShop,
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
		)

		orderHandler = orderH.NewOrderHandler(
			&c.FindOrders,
			&c.CreateOrder,
		)
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Get("/", chains.Core(productHandler.FindProducts))
			r.Post("/", chains.MerchantOnly(productHandler.CreateProduct))

			r.Route("/{slug}", func(r chi.Router) {
				r.Get("/", chains.Core(productHandler.GetProduct))
				r.Post("/images", chains.MerchantOnly(productHandler.AddProductImages))
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", chains.Core(authHandler.SignInEmail))
			r.Post("/merchant/signin", chains.Core(authHandler.SignInMerchantEmail))
			r.Post("/signup", chains.Core(authHandler.SignUpAccount))
			r.Post("/verify", chains.Core(authHandler.VerifyAccount))
			r.Post("/logout", chains.CoreAuth(authHandler.Logout))
			r.Get("/me", chains.CoreAuth(authHandler.Me))
		})

		r.Route("/merchants", func(r chi.Router) {
			r.Get("/", chains.MerchantAdminOnly(merchantHandler.FindMerchants))
			r.Post("/", chains.MerchantAdminOnly(merchantHandler.CreateMerchant))
			r.Route("/{merchantID}/accounts", func(r chi.Router) {
				r.Post("/", chains.MerchantAdminOnly(merchantHandler.AddMerchantAccount))
			})
		})

		r.Route("/customers", func(r chi.Router) {
			r.Get("/", chains.MerchantAdminOnly(customerHandler.FindCustomers))
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
				r.Delete("/{shopID}/{productID}", chains.CustomerOnly(cartHandler.RemoveItem))
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
			r.Get("/{id}", chains.MerchantAdminOnly(userHandler.GetUserByID))
		})

		r.Route("/users/me", func(r chi.Router) {
			r.Get("/", chains.CustomerOnly(userHandler.GetCurrentUser))

			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", chains.CustomerOnly(addressHandler.ListUserAddresses))
				r.Post("/", chains.CustomerOnly(addressHandler.SaveUserAddress))
				r.Delete("/{addressID}", chains.CustomerOnly(addressHandler.DeleteUserAddress))
			})
		})

		r.Route("/shops", func(r chi.Router) {
			r.Get("/", chains.Core(shopHandler.FindShops))
			r.Post("/", chains.MerchantOnly(shopHandler.SaveShop))

			r.Route("/{shopID}", func(r chi.Router) {
				r.Get("/", chains.Core(shopHandler.GetShopByID))

				r.Route("/addresses", func(r chi.Router) {
					r.Get("/", chains.Core(shopHandler.GetShopAddresses))
					r.Post("/", chains.MerchantOnly(addressHandler.CreateShopAddress))
				})

				r.Route("/couriers", func(r chi.Router) {
					r.Get("/", chains.Core(shopHandler.GetShopCouriers))
					r.Post("/", chains.MerchantOnly(courierHandler.ConfigureCourierShop))
				})

				r.Route("/products", func(r chi.Router) {
					r.Get("/", chains.Core(shopHandler.GetShopProducts))
					r.Post("/{productID}/inventories",
						chains.MerchantOnly(inventoryHandler.AddInventory))
				})
			})
		})

		r.Route("/midtrans", func(r chi.Router) {
			r.Post("/webhook", chains.Core(paymentHandler.HandleMidtransWebhook))
		})

		r.Route("/payments", func(r chi.Router) {
			r.Post("/{id}/action", chains.MerchantOnly(paymentHandler.ProcessManualPayment))

			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", chains.MerchantOnly(paymentHandler.ListPaymentAccount))
				r.Post("/", chains.MerchantAdminOnly(paymentHandler.CreatePaymentAccount))
			})

			r.Route("/methods", func(r chi.Router) {
				r.Get("/", chains.MerchantOnly(paymentHandler.ListPaymentMethod))
				r.Post("/", chains.MerchantAdminOnly(paymentHandler.CreatePaymentMethod))
			})
		})

		r.Route("/shipping", func(r chi.Router) {
			r.Post("/cost", chains.CoreAuth(shipmentHandler.EstimateShippingOptions))
		})

		r.Route("/order", func(r chi.Router) {
			r.Post("/", chains.CustomerOnly(orderHandler.CreateOrder))
		})

		r.Route("/orders", func(r chi.Router) {
			r.Get("/", chains.MerchantOnly(orderHandler.FindOrders))
		})
	})

	return r
}
