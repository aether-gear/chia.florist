package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	authendomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/user/usecase"

	"github.com/google/uuid"
)

type UserHandler struct {
	getUser       *usecase.GetUserUsecase
	findCustomers *usecase.FindCustomersUsecase
}

func NewUserHandler(
	getUser *usecase.GetUserUsecase,
	findCustomers *usecase.FindCustomersUsecase,
) *UserHandler {
	return &UserHandler{
		getUser:       getUser,
		findCustomers: findCustomers,
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	id, err := apphttp.ParamUUID(r, "id")
	if err != nil {
		return apperrors.NewBadRequest("invalid user id")
	}

	result, err := h.getUser.ByID(r.Context(), id)
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

	result, err := h.getUser.ByID(r.Context(), authCtx.UserID)
	if err != nil {
		return err
	}
	if result == nil {
		return apperrors.NewNotFound("user not found")
	}

	response := map[string]userResponse{
		"me": {
			ID:          result.ID,
			Name:        result.Name,
			Username:    result.Username,
			Phone:       result.Phone,
			LastLoginAt: result.LastLoginAt,
		},
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *UserHandler) FindCustomers(w http.ResponseWriter, r *http.Request) error {
	page := apphttp.QueryIntDefault(r, "page", 1)
	if page <= 0 {
		page = 1
	}
	limit := apphttp.QueryIntDefault(r, "limit", 10)
	if limit <= 0 {
		limit = 10
	}

	name := apphttp.Query(r, "name")
	username := apphttp.Query(r, "username")
	email := apphttp.Query(r, "email")
	idStr := apphttp.Query(r, "id")
	sort := apphttp.Query(r, "sort")

	input := usecase.FindCustomersInput{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
	if name != "" {
		input.Name = &name
	}
	if username != "" {
		input.Username = &username
	}
	if email != "" {
		input.Email = &email
	}
	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return apperrors.NewBadRequest("invalid user id")
		}
		input.ID = &id
	}

	users, total, err := h.findCustomers.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	results := make([]userResponse, 0, len(users))
	for _, u := range users {
		results = append(results, userResponse{
			ID:          u.ID,
			Name:        u.Name,
			Username:    u.Username,
			Phone:       u.Phone,
			LastLoginAt: u.LastLoginAt,
		})
	}

	response := map[string]interface{}{
		"users": results,
		"page":  page,
		"limit": limit,
		"total": total,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}
