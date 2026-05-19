package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
	"service-core/internal/common/middleware"

	addressH "service-core/internal/modules/address/delivery/http"
	authH "service-core/internal/modules/authentication/delivery/http"
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
		log  = c.Logger
		core = buildChain(log)
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
			r.Post("/", core(productHandler.CreateProduct))
			r.Get("/{id}", core(productHandler.GetProduct))
			r.Post("/images", core(productHandler.AddProductImages))
		})

		r.Route("/inventories", func(r chi.Router) {
			r.Post("/", core(inventoryHandler.CreateInventory))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/signin", core(authHandler.SignInEmail))
			r.Post("/signup", core(authHandler.SignUpAccount))
			r.Post("/verify", core(authHandler.VerifyAccount))
		})

		r.Route("/carts", func(r chi.Router) {
			r.Get("/", core(cartHandler.GetCart))
			r.Route("/items", func(r chi.Router) {
				r.Post("/", core(cartHandler.AddItem))
				r.Put("/", core(cartHandler.UpdateItem))
				r.Delete("/", core(cartHandler.RemoveItem))
			})
		})

		r.Route("/locations", func(r chi.Router) {
			r.Get("/provinces", core(locationHandler.Province))
			r.Get("/cities", core(locationHandler.City))
			r.Get("/districts", core(locationHandler.District))
			r.Get("/villages", core(locationHandler.Village))
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", core(userHandler.GetUserByID))
			r.Post("/addresses", core(addressHandler.CreateUserAddress))
		})

		r.Route("/shops", func(r chi.Router) {
			r.Post("/", core(shopHandler.CreateShop))
			r.Get("/{id}", core(shopHandler.GetShopByID))

			r.Route("/addresses", func(r chi.Router) {
				r.Get("/", core(addressHandler.ListShopAddresses))
				r.Post("/", core(addressHandler.CreateShopAddress))
				r.Get("/{id}", core(addressHandler.GetShopAddress))
			})
			r.Post("/couriers", core(courierHandler.ConfigureCourierShop))
		})

		r.Route("/payments", func(r chi.Router) {
			r.Route("/accounts", func(r chi.Router) {
				r.Get("/", core(paymentHandler.ListPaymentAccount))
				r.Post("/", core(paymentHandler.CreatePaymentAccount))
			})
			r.Route("/methods", func(r chi.Router) {
				r.Get("/", core(paymentHandler.ListPaymentMethod))
				r.Post("/", core(paymentHandler.CreatePaymentMethod))
			})
		})

		r.Route("/shipping", func(r chi.Router) {
			r.Post("/cost", core(shipmentHandler.EstimateShippingOptions))
		})

	})

	return r
}

func buildChain(log logger.Logger, extra ...middleware.Middleware) func(apphttp.AppHandler) http.HandlerFunc {
	base := []middleware.Middleware{
		middleware.Recovery(log),
		middleware.Logging(log),
		middleware.Response(),
	}

	mws := append(base, extra...)

	return func(h apphttp.AppHandler) http.HandlerFunc {
		return middleware.Chain(h, mws...)
	}
}
