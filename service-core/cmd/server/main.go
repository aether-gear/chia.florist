package main

import (
	"log"
	"net/http"

	"service-core/internal/modules/transaction/handler"
	"service-core/internal/modules/transaction/infra"
)

func main() {
	repo := infra.NewTransactionRepository()

	router := handler.NewRouter(repo)

	log.Println("service-core running on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
