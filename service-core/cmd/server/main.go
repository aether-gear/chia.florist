package main

import (
	"log"
	"net/http"

	productHandler "service-core/internal/features/product/delivery/http"
	productRepo "service-core/internal/features/product/infra/persistence"
	productUsecase "service-core/internal/features/product/usecase"
	database "service-core/internal/infra/db"

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

	productRepo := productRepo.NewProductRepository(db)
	findProducts := productUsecase.NewFindProductsUsecase(productRepo)
	getProduct := productUsecase.NewGetProductsUsecase(productRepo)
	productH := productHandler.NewProductHandler(findProducts, getProduct)

	transactionRepo := transactionInfra.NewTransactionRepository()
	transactionH := transactionHandler.NewTransactionHandler(transactionRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/products", productH.FindProducts)
	mux.HandleFunc("/products/", productH.GetProduct)

	mux.HandleFunc("/transactions", transactionH.GetAll)

	log.Println("service-core running on :8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
