package main

import (
	"log"
	"net/http"
	"time"

	aU "service-core/internal/features/auth/usecase"
	cU "service-core/internal/features/cart/usecase"
	pU "service-core/internal/features/product/usecase"
	"service-core/internal/shared/config"

	// uU "service-core/internal/features/user/usecase"

	aR "service-core/internal/features/auth/infra/persistence"
	cR "service-core/internal/features/cart/infra/persistence"
	pR "service-core/internal/features/product/infra/persistence"

	services "service-core/internal/features/auth/infra/service"

	database "service-core/internal/infra/db"

	authHandler "service-core/internal/features/auth/delivery/http"
	cartHandler "service-core/internal/features/cart/delivery/http"
	productHandler "service-core/internal/features/product/delivery/http"
	transactionHandler "service-core/internal/features/transaction/handler"

	transactionInfra "service-core/internal/features/transaction/infra"

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

	log.Println("service-core running on :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
