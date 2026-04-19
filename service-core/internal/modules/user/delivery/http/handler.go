package http

import (
	"net/http"
	"strings"

	"service-core/internal/modules/user/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

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

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		return errors.ErrBadRequest
	}

	id := parts[2]
	if id == "" {
		return errors.ErrBadRequest
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	result, err := h.getUser.ByID(parsedID)
	if err != nil {
		return err
	}
	if result == nil {
		return errors.ErrNotFound
	}

	response := UserResponse{
		ID:          result.ID,
		Name:        result.Name,
		Username:    result.Username,
		Phone:       result.Phone,
		LastLoginAt: result.LastLoginAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
