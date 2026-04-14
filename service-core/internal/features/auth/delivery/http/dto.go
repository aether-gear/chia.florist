package http

type SignUpParams struct {
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
}

type SignInEmailParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInUsernameParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
