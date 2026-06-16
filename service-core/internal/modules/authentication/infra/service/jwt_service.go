package service

import (
	"fmt"
	"strings"
	"time"

	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secret string) repository.TokenService {
	return &JWTService{
		secretKey: []byte(secret),
	}
}

type jwtClaims struct {
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	Type       string `json:"type"`
	MerchantID string `json:"merchant_id,omitempty"`
	Role       string `json:"roles,omitempty"` // comma-separated

	jwt.RegisteredClaims
}

func (j *JWTService) Generate(params repository.GenerateTokenParams) (repository.GeneratedToken, error) {
	now := time.Now()
	exp := now.Add(params.Duration)

	merchantIDStr := ""
	if params.MerchantID != nil {
		merchantIDStr = params.MerchantID.String()
	}

	roles := make([]string, len(params.Roles))
	for i, role := range params.Roles {
		roles[i] = string(role)
	}
	rolesStr := strings.Join(roles, ",")

	claims := jwtClaims{
		UserID:     params.UserID.String(),
		SessionID:  params.SessionID.String(),
		Type:       string(params.Type),
		MerchantID: merchantIDStr,
		Role:       rolesStr,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return repository.GeneratedToken{}, fmt.Errorf("sign jwt token failed: %w", err)
	}

	return repository.GeneratedToken{
		Token:     signed,
		ExpiresAt: exp,
		Type:      params.Type,
	}, nil
}

func (j *JWTService) Validate(tokenStr string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&jwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected jwt signing method")
			}

			return j.secretKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse jwt token failed: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, fmt.Errorf("invalid jwt claims type")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid jwt token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id claim: %w", err)
	}

	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session id claim: %w", err)
	}

	tknClaim := &domain.TokenClaims{
		UserID:    userID,
		SessionID: sessionID,
		Type:      domain.TokenType(claims.Type),
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	if claims.MerchantID != "" {
		mid, err := uuid.Parse(claims.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant_id claim: %w", err)
		}
		tknClaim.MerchantID = &mid
	}
	if claims.Role != "" {
		tknClaim.Roles = strings.Split(claims.Role, ",")
	}

	return tknClaim, nil
}
