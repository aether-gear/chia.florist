package http

import (
	"net/http"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	"service-core/internal/modules/customer/usecase"

	"github.com/google/uuid"
)

type CustomerHandler struct {
	findCustomers *usecase.FindCustomersUsecase
}

func NewCustomerHandler(
	findCustomers *usecase.FindCustomersUsecase,
) *CustomerHandler {
	return &CustomerHandler{
		findCustomers: findCustomers,
	}
}

func (h *CustomerHandler) FindCustomers(w http.ResponseWriter, r *http.Request) error {
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

	results := make([]customerResponse, 0, len(users))
	for _, u := range users {
		results = append(results, customerResponse{
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
