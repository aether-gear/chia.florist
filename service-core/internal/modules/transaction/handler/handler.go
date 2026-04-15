package handler

import (
	"encoding/json"
	"net/http"

	"service-core/internal/modules/transaction/application"
	"service-core/internal/modules/transaction/domain"
)

type TransactionHandler struct {
	usecase *application.TransactionUsecase
}

func NewTransactionHandler(repo domain.TransactionPort) *TransactionHandler {
	return &TransactionHandler{
		usecase: application.NewTransactionUsecase(repo),
	}
}

func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.usecase.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
