package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	"service-core/internal/common/logger"
	"service-core/internal/common/middleware"

	addressHandler "service-core/internal/modules/address/delivery/http"
	authHandler "service-core/internal/modules/auth/delivery/http"
	cartHandler "service-core/internal/modules/cart/delivery/http"
	locationHandler "service-core/internal/modules/location/delivery/http"
	paymentHandler "service-core/internal/modules/payment/delivery/http"
	productHandler "service-core/internal/modules/product/delivery/http"
	shopHandler "service-core/internal/modules/shop/delivery/http"
	userHandler "service-core/internal/modules/user/delivery/http"
)

func NewRouter(c *Container) *http.ServeMux {
	var (
		log  = c.Logger
		core = buildChain(log)
	)

	var (
		productH = productHandler.NewProductHandler(
			&c.FindProducts,
			&c.GetProduct,
			&c.CreateProduct,
		)

		authH = authHandler.NewAuthHandler(
			&c.LoginAccount,
			&c.RegisterAccount,
			&c.GetAccount,
		)

		cartH = cartHandler.NewCartHandler(
			&c.AddItem,
			&c.GetCart,
			&c.UpdateItem,
			&c.RemoveItem,
		)

		locationH = locationHandler.NewLocationHandler(
			&c.ListLocations,
		)

		userH = userHandler.NewUserHandler(
			&c.GetUser,
		)

		addressH = addressHandler.NewAddressHandler(
			&c.GetAddress,
			&c.GetShopAddress,
			&c.GetShopAddresses,
			&c.CreateAddress,
			&c.CreateShopAddress,
		)

		paymentH = paymentHandler.NewPaymentHandler(
			&c.CreatePaymentAccount,
			&c.ListPaymentAccount,
			&c.CreatePaymentMethod,
			&c.ListPaymentMethod,
		)

		shopH = shopHandler.NewAddressHandler(
			&c.CreateShop,
		)
	)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/product",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet:  productH.FindProducts,
			http.MethodPost: productH.CreateProduct,
		})),
	)
	mux.HandleFunc(
		"/product/",
		core(productH.GetProduct),
	)

	mux.HandleFunc(
		"/auth/signin",
		core(authH.SignInByEmail),
	)
	mux.HandleFunc(
		"/auth/signup",
		core(authH.SignUp),
	)

	mux.HandleFunc(
		"/cart",
		core(cartH.GetCart),
	)
	mux.HandleFunc(
		"/cart/items",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost:   cartH.AddItem,
			http.MethodPut:    cartH.UpdateItem,
			http.MethodDelete: cartH.RemoveItem,
		})),
	)

	mux.HandleFunc(
		"/locations/provinces",
		core(locationH.Province),
	)
	mux.HandleFunc(
		"/locations/cities/",
		core(locationH.City),
	)
	mux.HandleFunc(
		"/locations/districts/",
		core(locationH.District),
	)
	mux.HandleFunc(
		"/locations/villages/",
		core(locationH.Village),
	)

	mux.HandleFunc(
		"/user/",
		core(userH.GetUserByID),
	)
	mux.HandleFunc(
		"/user/address",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: addressH.CreateAddress,
		})),
	)

	mux.HandleFunc(
		"/shop",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: shopH.CreateShop,
		})),
	)
	mux.HandleFunc(
		"/shop/address",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodPost: addressH.CreateShopAddress,
		})),
	)
	mux.HandleFunc(
		"/shop/address/",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet: addressH.GetShopAddress,
		})),
	)
	mux.HandleFunc(
		"/shop/addresses/",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet: addressH.GetShopAddresses,
		})),
	)

	mux.HandleFunc(
		"/payment/account",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet:  paymentH.ListPaymentAccount,
			http.MethodPost: paymentH.CreatePaymentAccount,
		})),
	)
	mux.HandleFunc(
		"/payment/method",
		core(apphttp.HandleMethods(apphttp.MethodHandler{
			http.MethodGet:  paymentH.ListPaymentMethod,
			http.MethodPost: paymentH.CreatePaymentMethod,
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
