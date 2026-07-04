package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	apperrors "service-core/internal/common/errors"
	apphttp "service-core/internal/common/http"
	appcookie "service-core/internal/common/http/cookie"
	authdomain "service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/usecase"
	appconfig "service-core/internal/shared/config"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authHandler struct {
	me                *usecase.MeUsecase
	logout            *usecase.LogoutUsecase
	loginCustomer     *usecase.LoginCustomerUsecase
	loginStaff        *usecase.LoginStaffUsecase
	registerCustomer  *usecase.RegisterCustomerUsecase
	verifyAccount     *usecase.VerifyAccountUsecase
	getAccount        *usecase.GetAccountUsecase
	authenticateOAuth *usecase.AuthenticateOAuthUsecase
	googleCfg         appconfig.GoogleOAuthConfig
}

func NewAuthHandler(
	me *usecase.MeUsecase,
	logout *usecase.LogoutUsecase,
	loginCustomer *usecase.LoginCustomerUsecase,
	loginStaff *usecase.LoginStaffUsecase,
	registerCustomer *usecase.RegisterCustomerUsecase,
	verifyAccount *usecase.VerifyAccountUsecase,
	getAccount *usecase.GetAccountUsecase,
	authenticateOAuth *usecase.AuthenticateOAuthUsecase,
	googleCfg appconfig.GoogleOAuthConfig,
) *authHandler {
	return &authHandler{
		me:                me,
		logout:            logout,
		loginCustomer:     loginCustomer,
		loginStaff:        loginStaff,
		registerCustomer:  registerCustomer,
		verifyAccount:     verifyAccount,
		getAccount:        getAccount,
		authenticateOAuth: authenticateOAuth,
		googleCfg:         googleCfg,
	}
}

func (h *authHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok {
		return apperrors.NewUnauthorized("authentication required")
	}

	acc, err := h.getAccount.Execute(r.Context(), authCtx.UserID)
	if err != nil {
		return err
	}
	if acc == nil {
		return apperrors.NewNotFound("account not found")
	}

	response := map[string]any{
		"id":            acc.ID,
		"email":         acc.Email,
		"last_login_at": acc.LastLoginAt,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) Me(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	me, err := h.me.Execute(r.Context(), *authCtx)
	if err != nil {
		return err
	}

	roles := make([]roleResponse, 0, len(me.Actor.Roles))
	for _, role := range me.Actor.Roles {
		roles = append(roles, roleResponse{
			Code: string(role.Code),
			Name: role.Name,
		})
	}

	permissions := make([]permissionResponse, 0, len(me.Actor.Permissions))
	for _, role := range me.Actor.Roles {
		permissions = append(permissions, permissionResponse{
			Code: string(role.Code),
		})
	}

	response := meResponse{
		AccountID:       me.Account.ID,
		AccountType:     string(me.Account.Type),
		IsAuthenticated: true,
		StaffID:         me.Actor.StaffID,
		Roles:           roles,
		Permissions:     permissions,
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) SignInEmail(w http.ResponseWriter, r *http.Request) error {
	var req signInEmailRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("invalid email")
	}
	if req.Password == "" {
		return apperrors.NewBadRequest("invalid password")
	}

	input := usecase.LoginCustomerParams{
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
		Email:     req.Email,
		Password:  req.Password,
	}

	tokens, err := h.loginCustomer.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	// Access token is bound to the response cookie
	// so it can be used for authenticated API requests
	appcookie.Bind(
		w,
		appcookie.CookieAccess,
		tokens.AccessToken.Token,
		tokens.AccessToken.ExpiresAt,
	)

	// Refresh token is bound to the response cookie
	// so it can be used to obtain a new access token
	appcookie.Bind(
		w,
		appcookie.CookieCustomerRefresh,
		tokens.RefreshToken.Token,
		tokens.RefreshToken.ExpiresAt,
	)

	response := map[string]string{
		"message": "login success",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) SignUpAccount(w http.ResponseWriter, r *http.Request) error {
	var req signUpRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("invalid email")
	}
	if req.Password == "" {
		return apperrors.NewBadRequest("invalid password")
	}
	if req.Username == "" {
		return apperrors.NewBadRequest("invalid user name")
	}

	input := usecase.RegisterCustomerParams{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Username: req.Username,
		Phone:    req.Phone,
	}

	challengeID, err := h.registerCustomer.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	response := signUpResponse{
		Message:     "verification code sent",
		ChallengeID: *challengeID,
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *authHandler) VerifyAccount(w http.ResponseWriter, r *http.Request) error {
	var req verifyAccountRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.ChallengeID == "" {
		return apperrors.NewBadRequest("invalid challenge id")
	}
	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		return apperrors.NewBadRequest("invalid challenge id")
	}
	if len(req.OTP) != 6 {
		return apperrors.NewBadRequest("invalid otp")
	}

	input := usecase.VerifyAccountParams{
		UserAgent:   req.UserAgent,
		IPAddress:   req.IPAddress,
		ChallengeID: challengeID,
		OTP:         req.OTP,
	}

	tokens, err := h.verifyAccount.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	// Access token is bound to the response cookie
	// so it can be used for authenticated API requests
	appcookie.Bind(
		w,
		appcookie.CookieAccess,
		tokens.AccessToken.Token,
		tokens.AccessToken.ExpiresAt,
	)

	// Refresh token is bound to the response cookie
	// so it can be used to obtain a new access token
	appcookie.Bind(
		w,
		appcookie.CookieCustomerRefresh,
		tokens.RefreshToken.Token,
		tokens.RefreshToken.ExpiresAt,
	)

	response := map[string]string{
		"message": "verify success",
	}

	apphttp.WriteJSON(w, http.StatusCreated, response)
	return nil
}

func (h *authHandler) SignInStaffEmail(w http.ResponseWriter, r *http.Request) error {
	var req signInEmailRequest

	if err := apphttp.DecodeJSON(r, &req); err != nil {
		return apperrors.NewBadRequest("invalid body request")
	}

	if req.Email == "" {
		return apperrors.NewBadRequest("invalid email")
	}
	if req.Password == "" {
		return apperrors.NewBadRequest("invalid password")
	}

	input := usecase.LoginStaffParams{
		UserAgent: req.UserAgent,
		IPAddress: req.IPAddress,
		Email:     req.Email,
		Password:  req.Password,
	}

	tokens, err := h.loginStaff.Execute(r.Context(), input)
	if err != nil {
		return err
	}

	// Access token is bound to the response cookie
	// so it can be used for authenticated API requests
	appcookie.Bind(
		w,
		appcookie.CookieStaff,
		tokens.AccessToken.Token,
		tokens.AccessToken.ExpiresAt,
	)

	// Refresh token is bound to the response cookie
	// so it can be used to obtain a new access token
	appcookie.Bind(
		w,
		appcookie.CookieStaffRefresh,
		tokens.RefreshToken.Token,
		tokens.RefreshToken.ExpiresAt,
	)

	response := map[string]string{
		"message": "login success",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	authCtx, ok := authdomain.GetAuthContext(r.Context())
	if !ok || !authCtx.IsAuthenticated {
		return apperrors.NewUnauthorized("authentication required")
	}

	err := h.logout.Execute(r.Context(), *authCtx)
	if err != nil {
		return err
	}

	if authCtx.CustomerID != nil {
		appcookie.Clear(w, appcookie.CookieAccess)
		appcookie.Clear(w, appcookie.CookieCustomerRefresh)
	} else if authCtx.StaffID != nil {
		appcookie.Clear(w, appcookie.CookieStaff)
		appcookie.Clear(w, appcookie.CookieStaffRefresh)
	}

	response := map[string]string{
		"message": "logout success",
	}

	apphttp.WriteJSON(w, http.StatusOK, response)
	return nil
}

func (h *authHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) error {
	oauth2Config := &oauth2.Config{
		ClientID:     h.googleCfg.ClientID,
		ClientSecret: h.googleCfg.ClientSecret,
		RedirectURL:  h.googleCfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return apperrors.NewInternal(fmt.Errorf("failed to generate state: %w", err))
	}
	state := base64.URLEncoding.EncodeToString(b)

	appcookie.Bind(
		w,
		appcookie.CookieOAuthState,
		state,
		time.Now().Add(10*time.Minute),
	)

	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (h *authHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) error {
	cookie, err := appcookie.
		Extract(r, appcookie.CookieOAuthState)
	if err != nil {
		return apperrors.NewBadRequest("missing oauth state cookie")
	}

	stateParam := r.URL.Query().Get("state")
	if stateParam == "" || stateParam != cookie {
		return apperrors.NewBadRequest("invalid oauth state")
	}

	appcookie.Clear(w, appcookie.CookieOAuthState)

	code := r.URL.Query().Get("code")
	if code == "" {
		return apperrors.NewBadRequest("missing oauth code")
	}

	oauth2Config := &oauth2.Config{
		ClientID:     h.googleCfg.ClientID,
		ClientSecret: h.googleCfg.ClientSecret,
		RedirectURL:  h.googleCfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	token, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		return apperrors.NewUnauthorized(fmt.Sprintf("failed to exchange code: %v", err))
	}

	client := oauth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return apperrors.NewInternal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return apperrors.NewInternal(fmt.Errorf("google userinfo returned status code %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apperrors.NewInternal(err)
	}

	var googleUser struct {
		Sub     string  `json:"sub"`
		Name    string  `json:"name"`
		Email   string  `json:"email"`
		Picture *string `json:"picture"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return apperrors.NewInternal(err)
	}

	if googleUser.Email == "" {
		return apperrors.NewBadRequest("google did not provide email address")
	}

	userAgent := r.UserAgent()
	ipAddress := r.RemoteAddr

	params := usecase.AuthenticateOAuthParams{
		UserAgent: &userAgent,
		IPAddress: &ipAddress,
		Provider:  authdomain.OAuthProviderGoogle,
		Subject:   googleUser.Sub,
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: googleUser.Picture,
	}

	result, err := h.authenticateOAuth.
		Execute(r.Context(), params)
	if err != nil {
		return err
	}

	// Access token is bound to the response cookie
	// so it can be used for authenticated API requests
	appcookie.Bind(
		w,
		appcookie.CookieAccess,
		result.AccessToken.Token,
		result.AccessToken.ExpiresAt,
	)

	// Refresh token is bound to the response cookie
	// so it can be used to obtain a new access token
	appcookie.Bind(
		w,
		appcookie.CookieCustomerRefresh,
		result.RefreshToken.Token,
		result.RefreshToken.ExpiresAt,
	)

	successRedirectURL := h.googleCfg.SuccessRedirectURL
	if successRedirectURL == "" {
		successRedirectURL = "/"
	}
	http.Redirect(w, r, successRedirectURL, http.StatusTemporaryRedirect)
	return nil
}
