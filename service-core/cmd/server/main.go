package main

import (
	"log"
	"net/http"
	"time"

	adU "service-core/internal/modules/address/usecase"
	aU "service-core/internal/modules/auth/usecase"
	cU "service-core/internal/modules/cart/usecase"
	lU "service-core/internal/modules/location/usecase"
	pU "service-core/internal/modules/product/usecase"
	uU "service-core/internal/modules/user/usecase"
	"service-core/internal/shared/config"

	// uU "service-core/internal/modules/user/usecase"

	adR "service-core/internal/modules/address/infra/persistence"
	aR "service-core/internal/modules/auth/infra/persistence"
	cR "service-core/internal/modules/cart/infra/persistence"
	lR "service-core/internal/modules/location/infra/persistence"
	pR "service-core/internal/modules/product/infra/persistence"
	uR "service-core/internal/modules/user/infra/persistence"

	services "service-core/internal/modules/auth/infra/service"

	database "service-core/internal/infra/db"

	addressHandler "service-core/internal/modules/address/delivery/http"
	authHandler "service-core/internal/modules/auth/delivery/http"
	cartHandler "service-core/internal/modules/cart/delivery/http"
	locationHandler "service-core/internal/modules/location/delivery/http"
	productHandler "service-core/internal/modules/product/delivery/http"
	transactionHandler "service-core/internal/modules/transaction/handler"
	userHandler "service-core/internal/modules/user/delivery/http"

	transactionInfra "service-core/internal/modules/transaction/infra"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}

	cfg := database.LoadConfig()

	db := database.NewConnection(cfg)
	defer db.Close()

	productRepo := pR.NewProductRepository(db)
	authRepo := aR.NewAuthRepository(db)
	transactionRepo := transactionInfra.NewTransactionRepository()
	cartRepo := cR.NewCartRepositoryImpl(db)
	locationRepo := lR.NewLocationRepositoryImpl(db)
	userRepo := uR.NewUserRepositoryImpl(db)
	addressRepo := adR.NewAddressRepositoryImpl(db)

	tokenSvc := services.NewJWTService(
		config.MustGetEnv("JWT_SECRET"),
		24*time.Minute,
	)

	hasher := services.NewBcryptHasher()

	findProducts := pU.NewFindProductsUsecase(productRepo)
	getProduct := pU.NewGetProductsUsecase(productRepo)
	createProdoct := pU.NewCreateProductUsecase(productRepo)
	loginAccount := aU.NewLoginUsecase(
		authRepo,
		hasher,
		tokenSvc,
	)
	registerAccount := aU.NewRegisterUsecase(
		authRepo,
		hasher,
	)
	getAccount := aU.NewGetAccountUsecase(
		authRepo,
	)
	getItemCart := cU.NewGetCartUsecase(cartRepo, productRepo)
	addItemCart := cU.NewAddItemUsecase(cartRepo, productRepo)
	updateItemCart := cU.NewUpdateItemUsecase(cartRepo, productRepo)
	removeItemCart := cU.NewRemoveItemUsecase(cartRepo)
	listLocations := lU.NewListLocationUsecase(locationRepo)
	getUser := uU.NewGetUserUsecase(userRepo)
	getAddress := adU.NewGetAddressUsecase(addressRepo)
	createAddress := adU.NewCreateAddressUsecase(addressRepo)

	// findUsers := uU.NewUser

	productH := productHandler.NewProductHandler(
		findProducts,
		getProduct,
		createProdoct,
	)
	transactionH := transactionHandler.NewTransactionHandler(transactionRepo)
	authH := authHandler.NewAuthHandler(
		loginAccount,
		registerAccount,
		getAccount,
	)
	cartH := cartHandler.NewCartHandler(addItemCart, getItemCart, updateItemCart, removeItemCart)
	locationH := locationHandler.NewLocationHandler(listLocations)
	userH := userHandler.NewUserHandler(getUser)
	addressH := addressHandler.NewAddressHandler(
		getAddress,
		createAddress,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			productH.FindProducts(w, r)
		}
		if r.Method == http.MethodPost {
			productH.CreateProduct(w, r)
		}
	})
	mux.HandleFunc("/product/", productH.GetProduct)

	mux.HandleFunc("/transactions", transactionH.GetAll)

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

	log.Println("service-core running on :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
