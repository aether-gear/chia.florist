package http

import (
	"encoding/json"
	"net/http"

	"service-core/internal/modules/auth/usecase"

	"service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"

	"github.com/google/uuid"
)

type authHandler struct {
	signInEmail *usecase.LoginEmailUsecase
	signUp      *usecase.RegisterUsecase
	getAccount  *usecase.GetAccountUsecase
}

func NewAuthHandler(
	signInEmail *usecase.LoginEmailUsecase,
	signUp *usecase.RegisterUsecase,
	getAccount *usecase.GetAccountUsecase,
) *authHandler {
	return &authHandler{
		signInEmail: signInEmail,
		signUp:      signUp,
		getAccount:  getAccount,
	}
}

func (h *authHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	UserId, ok := r.Context().Value("user_id").(string)
	if !ok {
		return errors.ErrUnauthorized
	}

	userID, err := uuid.Parse(UserId)
	if err != nil {
		return errors.ErrBadRequest
	}

	acc, err := h.getAccount.Execute(userID)
	if err != nil {
		return err
	}

	if acc == nil {
		return errors.ErrNotFound
	}

	response := map[string]interface{}{
		"id":            acc.ID,
		"email":         acc.Email,
		"last_login_at": acc.LastLoginAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) SignInEmail(w http.ResponseWriter, r *http.Request) error {
	var req SignInEmailParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" {
		return errors.ErrBadRequest
	}

	token, exp, err := h.signInEmail.Execute(req.Email, req.Password)
	if err != nil {
		return errors.ErrUnauthorized
	}

	response := map[string]interface{}{
		"access_token": token,
		"expires_in":   exp,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) SignUpAccount(w http.ResponseWriter, r *http.Request) error {
	var req SignUpParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.ErrBadRequest
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		return errors.ErrBadRequest
	}

	err := h.signUp.Execute(usecase.SignUpParams{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
	})
	if err != nil {
		return err
	}

	response := map[string]string{
		"message": "account successfully created",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}
