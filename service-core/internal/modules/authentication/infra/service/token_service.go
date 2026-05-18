package service

import (
	"fmt"
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

func (j *JWTService) Generate(params repository.GenerateTokenParams) (repository.GeneratedToken, error) {
	now := time.Now()
	exp := now.Add(params.Duration)

	claims := domain.TokenClaims{
		UserID:    params.UserID.String(),
		SessionID: params.SessionID.String(),
		Type:      params.Type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
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

	result := repository.GeneratedToken{
		Token:     signed,
		ExpiresAt: exp,
		Type:      params.Type,
	}

	return result, nil
}

func (j *JWTService) Validate(tokenStr string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&domain.TokenClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse jwt token failed: %w", err)
	}

	claims, ok := token.Claims.(*domain.TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid jwt token claims")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in claims: %w", err)
	}

	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session id in claims: %w", err)
	}

	return &domain.TokenClaims{
		UserID:    userID.String(),
		SessionID: sessionID.String(),
	}, nil
}
