package main

import (
	"log"
	"net/http"
	"time"

	aU "service-core/internal/features/auth/usecase"
	pU "service-core/internal/features/product/usecase"
	"service-core/internal/shared/config"

	// uU "service-core/internal/features/user/usecase"

	aR "service-core/internal/features/auth/infra/persistence"
	pR "service-core/internal/features/product/infra/persistence"

	services "service-core/internal/features/auth/infra/service"

	database "service-core/internal/infra/db"

	authHandler "service-core/internal/features/auth/delivery/http"
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

	tokenSvc := services.NewJWTService(
		config.MustGetEnv("JWT_SECRET"),
		24*time.Minute,
	)

	hasher := services.NewBcryptHasher()

	findProducts := pU.NewFindProductsUsecase(productRepo)
	getProduct := pU.NewGetProductsUsecase(productRepo)
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

	// findUsers := uU.NewUser

	productH := productHandler.NewProductHandler(findProducts, getProduct)
	transactionH := transactionHandler.NewTransactionHandler(transactionRepo)
	authH := authHandler.NewAuthHandler(
		loginAccount,
		registerAccount,
		getAccount,
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/products", productH.FindProducts)
	mux.HandleFunc("/products/", productH.GetProduct)

	mux.HandleFunc("/transactions", transactionH.GetAll)

	mux.HandleFunc("/auth/signin", authH.SignInByEmail)
	mux.HandleFunc("/auth/signup", authH.SignUp)

	log.Println("service-core running on :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
