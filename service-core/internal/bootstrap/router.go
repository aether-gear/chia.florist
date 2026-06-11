package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	applogger "service-core/internal/common/logger"
	appmiddleware "service-core/internal/common/middleware"

	authendomain "service-core/internal/modules/authentication/domain"
	authorzDomain "service-core/internal/modules/authorization/domain"

	addressH "service-core/internal/modules/address/delivery/http"
	authH "service-core/internal/modules/authentication/delivery/http"
	cartH "service-core/internal/modules/cart/delivery/http"
	courierH "service-core/internal/modules/courier/delivery/http"
	inventoryH "service-core/internal/modules/inventory/delivery/http"
	locationH "service-core/internal/modules/location/delivery/http"
	merchantH "service-core/internal/modules/merchant/delivery/http"
	paymentH "service-core/internal/modules/payment/delivery/http"
	productH "service-core/internal/modules/product/delivery/http"
	shipmentH "service-core/internal/modules/shipment/delivery/http"
	shopH "service-core/internal/modules/shop/delivery/http"
	userH "service-core/internal/modules/user/delivery/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(c *Container) *chi.Mux {
	var log = c.Logger

	var (
		core = buildChain(
			log,
			c.CORSAllowedOrigins,
		)
		coreAuth = buildChain(
			log,
			c.CORSAllowedOrigins,
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
		)
		merchantOnly = buildChain(
			log,
			c.CORSAllowedOrigins,
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				appcookie.AccessTokenCookieName,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeMerchant),
		)
		merchantAdminOnly = buildChain(
			log,
			c.CORSAllowedOrigins,
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				appcookie.AccessTokenMerchantCookieName,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeMerchant),
			c.Authorizer.RequireMerchantRole(authorzDomain.RoleMerchantAdmin),
		)
		customerOnly = buildChain(
			log,
			c.CORSAllowedOrigins,
			c.Authenticator.RequireAuth(
				c.DBExecutor,
				appcookie.AccessTokenCookieName,
			),
			c.Authorizer.LoadActor(c.DBExecutor),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeCustomer),
		)
	)

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
			&c.LoginCustomer,
			&c.LoginMerchant,
			&c.RegisterCustomer,
			&c.VerifyAccount,
			&c.GetAccount,
		)

		merchantHandler = merchantH.NewMerchantHandler(
			&c.AddMerchantAccount,
			&c.CreateMerchant,
		)

		cartHandler = cartH.NewCartHandler(
			&c.AddItem,
			&c.GetCart,
			&c.UpdateItem,
			&c.RemoveItem,
		)

		locationHandler = locationH.NewLocationHandler(
			&c.ListLocations,
		)

		userHandler = userH.NewUserHandler(
			&c.GetUser,
		)

		addressHandler = addressH.NewAddressHandler(
			&c.GetShopAddress,
			&c.ListUserAddresses,
			&c.ListShopAddresses,
			&c.CreateAddress,
			&c.CreateShopAddress,
		)

		paymentHandler = paymentH.NewPaymentHandler(
			&c.CreatePaymentAccount,
			&c.ListPaymentAccount,
			&c.CreatePaymentMethod,
			&c.ListPaymentMethod,
		)

		shopHandler = shopH.NewAddressHandler(
			&c.GetShop,
			&c.CreateShop,
		)

		courierHandler = courierH.NewCourierHandler(
			&c.ConfigureShopCourier,
		)

		shipmentHandler = shipmentH.NewShipmentHandler(
			&c.EstimateShippingOptions,
		)
	)

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Route("/products", func(r chi.Router) {
			r.Get("/", core(productHandler.FindProducts))
			r.Post("/", merchantOnly(productHandler.CreateProduct))
			r.Get("/{id}", core(productHandler.GetProduct))

			r.Post("/{id}/images", merchantOnly(productHandler.AddProductImages))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Get("/me", coreAuth(authHandler.Me))
			r.Post("/signin", core(authHandler.SignInEmail))
			r.Post("/merchant/signin", core(authHandler.SignInMerchantEmail))
			r.Post("/signup", core(authHandler.SignUpAccount))
			r.Post("/verify", core(authHandler.VerifyAccount))
		})

		r.Route("/merchants", func(r chi.Router) {
			r.Post("/", merchantAdminOnly(merchantHandler.CreateMerchant))
			r.Route("/{merchantID}/accounts", func(r chi.Router) {
				r.Post("/", merchantAdminOnly(merchantHandler.AddMerchantAccount))
			})
		})

		r.Route("/carts", func(r chi.Router) {
			r.Get("/", customerOnly(cartHandler.GetCart))

			r.Route("/items", func(r chi.Router) {
				r.Post("/", customerOnly(cartHandler.AddItem))
				r.Put("/{shopID}/{productID}", customerOnly(cartHandler.UpdateItem))
				r.Delete("/{shopID}/{productID}", customerOnly(cartHandler.RemoveItem))
			})
		})

		r.Route("/provinces", func(r chi.Router) {
			r.Get("/", core(locationHandler.Province))
			r.Get("/{id}/cities", core(locationHandler.City))
		})
		r.Route("/cities", func(r chi.Router) {
			r.Get("/{id}/districts", core(locationHandler.District))
		})
		r.Route("/districts", func(r chi.Router) {
			r.Get("/{id}/villages", core(locationHandler.Village))
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", merchantOnly(userHandler.GetUserByID))
		})

		r.Route("/users/me", func(r chi.Router) {
			r.Get("/", customerOnly(userHandler.GetCurrentUser))

			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", customerOnly(addressHandler.ListUserAddresses))
				r.Post("/", customerOnly(addressHandler.CreateUserAddress))
			})
		})

		r.Route("/shops", func(r chi.Router) {
			r.Post("/", merchantOnly(shopHandler.CreateShop))
			r.Get("/{id}", core(shopHandler.GetShopByID))

			r.Route("/{id}/addresses", func(r chi.Router) {
				r.Get("/", core(addressHandler.ListShopAddresses))
				r.Post("/", merchantOnly(addressHandler.CreateShopAddress))
				r.Get("/{addressID}", core(addressHandler.GetShopAddress))
			})

			r.Route("/{id}/couriers", func(r chi.Router) {
				r.Post("/", merchantOnly(courierHandler.ConfigureCourierShop))
			})

			r.Route("/{id}/products", func(r chi.Router) {
				r.Post("/{productID}/inventories",
					merchantOnly(inventoryHandler.AddInventory))
			})
		})

		r.Route("/payments", func(r chi.Router) {
			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", customerOnly(paymentHandler.ListPaymentAccount))
				r.Post("/", merchantOnly(paymentHandler.CreatePaymentAccount))
			})

			r.Route("/methods", func(r chi.Router) {
				r.Get("/", core(paymentHandler.ListPaymentMethod))
				r.Post("/", merchantOnly(paymentHandler.CreatePaymentMethod))
			})
		})

		r.Route("/shipping", func(r chi.Router) {
			r.Post("/cost", customerOnly(shipmentHandler.EstimateShippingOptions))
		})
	})

	return r
}

func buildChain(
	log applogger.Logger,
	allowedOrigins []string,
	extra ...appmiddleware.Middleware,
) func(apphttp.AppHandler) http.HandlerFunc {
	base := []appmiddleware.Middleware{
		appmiddleware.CORS(allowedOrigins),
		appmiddleware.Recovery(log),
		appmiddleware.Logging(log),
		appmiddleware.Response(),
	}

	mws := append(base, extra...)

	return func(h apphttp.AppHandler) http.HandlerFunc {
		return appmiddleware.Chain(h, mws...)
	}
}
