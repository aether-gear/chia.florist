package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"service-core/internal/modules/auth/usecase"

	"github.com/google/uuid"
)

type authHandler struct {
	signIn *usecase.LoginUsecase
	signUp *usecase.RegisterUsecase
	me     *usecase.GetAccountUsecase
}

func NewAuthHandler(
	signIn *usecase.LoginUsecase,
	signUp *usecase.RegisterUsecase,
	me *usecase.GetAccountUsecase,
) *authHandler {
	return &authHandler{
		signIn: signIn,
		signUp: signUp,
		me:     me,
	}
}

func (h *authHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	UserId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	parsedID, err := uuid.Parse(UserId)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	acc, err := h.me.ById(parsedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if acc == nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"id":            acc.ID,
		"email":         acc.Email,
		"last_login_at": acc.LastLoginAt,
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *authHandler) SignInByEmail(w http.ResponseWriter, r *http.Request) {
	var req SignInEmailParams

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	token, exp, err := h.signIn.ByEmail(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	resp := map[string]interface{}{
		"access_token": token,
		"expires_in":   exp,
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpParams

	body, _ := io.ReadAll(r.Body)

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	err := h.signUp.Register(usecase.SignUpParams{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "account successfully created",
	})
}
