package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
	"service-core/internal/common/middleware"

	addressH "service-core/internal/modules/address/delivery/http"
	authH "service-core/internal/modules/authentication/delivery/http"
	authendomain "service-core/internal/modules/authentication/domain"

	cartH "service-core/internal/modules/cart/delivery/http"
	courierH "service-core/internal/modules/courier/delivery/http"
	inventoryH "service-core/internal/modules/inventory/delivery/http"
	locationH "service-core/internal/modules/location/delivery/http"
	paymentH "service-core/internal/modules/payment/delivery/http"
	productH "service-core/internal/modules/product/delivery/http"
	shipmentH "service-core/internal/modules/shipment/delivery/http"
	shopH "service-core/internal/modules/shop/delivery/http"
	userH "service-core/internal/modules/user/delivery/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(c *Container) *chi.Mux {
	var (
		log          = c.Logger
		core         = buildChain(log, c.CORSAllowedOrigins)
		merchantOnly = buildChain(
			log,
			c.CORSAllowedOrigins,
			c.AuthMiddleware,
			c.Authorizer.LoadActor(),
			c.Authorizer.RequireAccountType(authendomain.AccountTypeMerchant),
		)
		customerOnly = buildChain(
			log,
			c.CORSAllowedOrigins,
			c.AuthMiddleware,
			c.Authorizer.LoadActor(),
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
			&c.LoginAccount,
			&c.RegisterAccount,
			&c.VerifyAccount,
			&c.GetAccount,
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
			r.Post("/signin", core(authHandler.SignInEmail))
			r.Post("/signup", core(authHandler.SignUpAccount))
			r.Post("/verify", core(authHandler.VerifyAccount))
		})

		r.Route("/carts", func(r chi.Router) {
			r.Get("/", customerOnly(cartHandler.GetCart))

			r.Route("/items", func(r chi.Router) {
				r.Post("/", customerOnly(cartHandler.AddItem))
				r.Put("/{itemID}", customerOnly(cartHandler.UpdateItem))
				r.Delete("/{itemID}", customerOnly(cartHandler.RemoveItem))
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
			r.Get("/{id}", core(userHandler.GetUserByID))
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
				r.Get("/", core(paymentHandler.ListPaymentAccount))
				r.Post("/", merchantOnly(paymentHandler.CreatePaymentAccount))
			})

			r.Route("/methods", func(r chi.Router) {
				r.Get("/", core(paymentHandler.ListPaymentMethod))
				r.Post("/", merchantOnly(paymentHandler.CreatePaymentMethod))
			})
		})

		r.Route("/shipping", func(r chi.Router) {
			r.Post("/cost", core(shipmentHandler.EstimateShippingOptions))
		})
	})

	return r
}

func buildChain(
	log logger.Logger,
	allowedOrigins []string,
	extra ...middleware.Middleware,
) func(apphttp.AppHandler) http.HandlerFunc {
	base := []middleware.Middleware{
		middleware.CORS(allowedOrigins),
		middleware.Recovery(log),
		middleware.Logging(log),
		middleware.Response(),
	}

	mws := append(base, extra...)

	return func(h apphttp.AppHandler) http.HandlerFunc {
		return middleware.Chain(h, mws...)
	}
}
