package main

import (
	"log"
	"net/http"

	productUsecase "service-core/internal/features/product/application"
	productHandler "service-core/internal/features/product/handler/http"
	productInfra "service-core/internal/features/product/infra"

	transactionHandler "service-core/internal/features/transaction/handler"
	transactionInfra "service-core/internal/features/transaction/infra"
)

func main() {
	productRepo := productInfra.NewProductRepository()
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
