package handler

import (
	"net/http"
	"service-core/internal/modules/transaction/domain"

	"github.com/go-chi/chi/v5"
)

func NewRouter(repo domain.TransactionPort) http.Handler {
	r := chi.NewRouter()

	handler := NewTransactionHandler(repo)
	r.Get("/transactions", handler.GetAll)

	return r
}
