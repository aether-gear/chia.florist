package bootstrap

import (
	"net/http"

	apphttp "service-core/internal/common/http"
	"service-core/internal/common/middleware"
	addressHandler "service-core/internal/modules/address/delivery/http"
	authHandler "service-core/internal/modules/auth/delivery/http"
	cartHandler "service-core/internal/modules/cart/delivery/http"
	locationHandler "service-core/internal/modules/location/delivery/http"
	productHandler "service-core/internal/modules/product/delivery/http"
	userHandler "service-core/internal/modules/user/delivery/http"
)

func NewRouter(c *Container) *http.ServeMux {
	mux := http.NewServeMux()

	productH := productHandler.NewProductHandler(
		&c.FindProducts,
		&c.GetProduct,
		&c.CreateProduct,
	)

	authH := authHandler.NewAuthHandler(
		&c.LoginAccount,
		&c.RegisterAccount,
		&c.GetAccount,
	)

	cartH := cartHandler.NewCartHandler(
		&c.AddItem,
		&c.GetCart,
		&c.UpdateItem,
		&c.RemoveItem,
	)

	locationH := locationHandler.NewLocationHandler(
		&c.ListLocations,
	)

	userH := userHandler.NewUserHandler(
		&c.GetUser,
	)

	addressH := addressHandler.NewAddressHandler(
		&c.GetAddress,
		&c.CreateAddress,
	)

	chain := middleware.Chain
	log := c.Logger

	mux.HandleFunc(
		"/product",
		chain(
			apphttp.HandleMethods(apphttp.MethodHandler{
				http.MethodGet:  productH.FindProducts,
				http.MethodPost: productH.CreateProduct,
			}),
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/product/",
		chain(
			productH.GetProduct,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/auth/signin",
		chain(
			authH.SignInByEmail,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/auth/signup",
		chain(
			authH.SignUp,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/cart",
		chain(
			cartH.GetCart,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/cart/items",
		chain(
			apphttp.HandleMethods(apphttp.MethodHandler{
				http.MethodPost:   cartH.AddItem,
				http.MethodPut:    cartH.UpdateItem,
				http.MethodDelete: cartH.RemoveItem,
			}),
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc("/locations/provinces",
		chain(locationH.Province,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc("/locations/cities/",
		chain(locationH.City,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc("/locations/districts/",
		chain(locationH.District,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc("/locations/villages/",
		chain(locationH.Village,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/user/",
		chain(
			userH.GetUserByID,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/user/addresses/",
		chain(
			addressH.GetAddresses,
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	mux.HandleFunc(
		"/user/address",
		chain(
			apphttp.HandleMethods(apphttp.MethodHandler{
				http.MethodPost: addressH.CreateAddress,
				// future todo:
				// http.MethodPut: addressH.UpdateAddress,
				// http.MethodDelete: addressH.DeleteAddress,
			}),
			middleware.Recovery(log),
			middleware.Logging(log),
			middleware.Response(),
		),
	)

	return mux
}
