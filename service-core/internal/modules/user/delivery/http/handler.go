package http

import (
	"encoding/json"
	"net/http"
	"service-core/internal/modules/user/usecase"
	"strings"

	"github.com/google/uuid"
)

type UserHandler struct {
	getUser *usecase.GetUserUsecase
}

func NewUserHandler(getUser *usecase.GetUserUsecase) *UserHandler {
	return &UserHandler{
		getUser: getUser,
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id := parts[2]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	result, err := h.getUser.ByID(parsedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	response := UserResponse{
		ID:          result.ID,
		Name:        result.Name,
		Username:    result.Username,
		Phone:       result.Phone,
		LastLoginAt: result.LastLoginAt,
	}

	json.NewEncoder(w).Encode(response)
}
