package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/user/usecase"
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
	id, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid user id")
	}

	result, err := h.getUser.ByID(id)
	if err != nil {
		return err
	}
	if result == nil {
		return apperrors.NewNotFound("user not found")
	}

	response := userResponse{
		ID:          result.ID,
		Name:        result.Name,
		Username:    result.Username,
		Phone:       result.Phone,
		LastLoginAt: result.LastLoginAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authendomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	result, err := h.getUser.ByID(authCtx.UserID)
	if err != nil {
		return err
	}
	if result == nil {
		return apperrors.NewNotFound("user not found")
	}

	response := userResponse{
		ID:          result.ID,
		Name:        result.Name,
		Username:    result.Username,
		Phone:       result.Phone,
		LastLoginAt: result.LastLoginAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
