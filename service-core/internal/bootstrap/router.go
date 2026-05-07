package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
	"service-core/internal/common/middleware"

	addressH "service-core/internal/modules/address/delivery/http"
	authH "service-core/internal/modules/auth/delivery/http"
	cartH "service-core/internal/modules/cart/delivery/http"
	courierH "service-core/internal/modules/courier/delivery/http"
	inventoryH "service-core/internal/modules/inventory/delivery/http"
	locationH "service-core/internal/modules/location/delivery/http"
	paymentH "service-core/internal/modules/payment/delivery/http"
	productH "service-core/internal/modules/product/delivery/http"
	shipmentH "service-core/internal/modules/shipment/delivery/http"
	shopH "service-core/internal/modules/shop/delivery/http"
	userH "service-core/internal/modules/user/delivery/http"
)

func NewRouter(c *Container) *http.ServeMux {
	var (
		log  = c.Logger
		core = buildChain(log)
	)

	var (
		productHandler = productH.NewProductHandler(
			&c.FindProducts,
			&c.GetProduct,
			&c.CreateProduct,
		)

		inventoryHandler = inventoryH.NewInventoryHandler(
			&c.CreateInventory,
		)

		authHandler = authH.NewAuthHandler(
			&c.LoginAccount,
			&c.RegisterAccount,
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

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/product",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet:  productHandler.FindProducts,
			http.MethodPost: productHandler.CreateProduct,
		})),
	)
	mux.HandleFunc(
		"/product/",
		core(productHandler.GetProduct),
	)
	mux.HandleFunc(
		"/inventory",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: inventoryHandler.CreateInventory,
		})),
	)

	mux.HandleFunc(
		"/auth/signin",
		core(authHandler.SignInEmail),
	)
	mux.HandleFunc(
		"/auth/signup",
		core(authHandler.SignUpAccount),
	)

	mux.HandleFunc(
		"/cart",
		core(cartHandler.GetCart),
	)
	mux.HandleFunc(
		"/cart/items",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost:   cartHandler.AddItem,
			http.MethodPut:    cartHandler.UpdateItem,
			http.MethodDelete: cartHandler.RemoveItem,
		})),
	)

	mux.HandleFunc(
		"/locations/provinces",
		core(locationHandler.Province),
	)
	mux.HandleFunc(
		"/locations/cities/",
		core(locationHandler.City),
	)
	mux.HandleFunc(
		"/locations/districts/",
		core(locationHandler.District),
	)
	mux.HandleFunc(
		"/locations/villages/",
		core(locationHandler.Village),
	)

	mux.HandleFunc(
		"/user/",
		core(userHandler.GetUserByID),
	)
	mux.HandleFunc(
		"/user/address",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: addressHandler.CreateUserAddress,
		})),
	)

	mux.HandleFunc(
		"/shop",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: shopHandler.CreateShop,
		})),
	)
	mux.HandleFunc(
		"/shop/",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet: shopHandler.GetShopByID,
		})),
	)
	mux.HandleFunc(
		"/shop/address",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: addressHandler.CreateShopAddress,
		})),
	)
	mux.HandleFunc(
		"/shop/address/",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet: addressHandler.GetShopAddress,
		})),
	)
	mux.HandleFunc(
		"/shop/addresses/",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet: addressHandler.ListShopAddresses,
		})),
	)

	mux.HandleFunc(
		"/payment/account",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet:  paymentHandler.ListPaymentAccount,
			http.MethodPost: paymentHandler.CreatePaymentAccount,
		})),
	)
	mux.HandleFunc(
		"/payment/method",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet:  paymentHandler.ListPaymentMethod,
			http.MethodPost: paymentHandler.CreatePaymentMethod,
		})),
	)

	mux.HandleFunc(
		"/shops/couriers",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: courierHandler.ConfigureCourierShop,
		})),
	)

	mux.HandleFunc(
		"/shipping/cost",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: shipmentHandler.EstimateShippingOptions,
		})),
	)

	return mux
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
