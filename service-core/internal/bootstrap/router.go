package bootstrap

import (
	"net/http"

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

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			productH.FindProducts(w, r)
		}
		if r.Method == http.MethodPost {
			productH.CreateProduct(w, r)
		}
	})
	mux.HandleFunc("/product/", productH.GetProduct)

	mux.HandleFunc("/auth/signin", authH.SignInByEmail)
	mux.HandleFunc("/auth/signup", authH.SignUp)

	mux.HandleFunc("/cart", cartH.GetCart)
	mux.HandleFunc("/cart/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			cartH.AddItem(w, r)
		case http.MethodPut:
			cartH.UpdateItem(w, r)
		case http.MethodDelete:
			cartH.RemoveItem(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/locations/provinces", locationH.Province)
	mux.HandleFunc("/locations/cities/", locationH.City)
	mux.HandleFunc("/locations/districts/", locationH.District)
	mux.HandleFunc("/locations/villages/", locationH.Village)

	mux.HandleFunc("/user/", userH.GetUserByID)
	mux.HandleFunc("/user/addresses/", addressH.GetAddresses)
	mux.HandleFunc("/user/address", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addressH.CreateAddress(w, r)
		case http.MethodPut:
		case http.MethodDelete:
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
