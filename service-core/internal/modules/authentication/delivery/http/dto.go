package http

import "github.com/google/uuid"

type SignUpParams struct {
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Phone    *string `json:"phone"`
}

type VerifyParams struct {
	UserAgent   *string `json:"user_agent"`
	IPAddress   *string `json:"ip_address"`
	ChallengeID string  `json:"challenge_id"`
	OTP         int     `json:"otp"`
}

type SignInEmailParams struct {
	UserAgent *string `json:"user_agent"`
	IPAddress *string `json:"ip_address"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
}

type SignInUsernameParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Message     string    `json:"message"`
	ChallengeID uuid.UUID `json:"challenge_id"`
}
