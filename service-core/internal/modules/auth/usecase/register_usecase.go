package usecase

import (
	"fmt"
	"time"

	authDomain "service-core/internal/modules/auth/domain"
	authRepo "service-core/internal/modules/auth/repository"

	"github.com/google/uuid"
)

type RegisterUsecase struct {
	authRepo authRepo.AuthRepository
	hasher   authDomain.PasswordHasher
}

func NewRegisterUsecase(
	authRepo authRepo.AuthRepository,
	hasher authDomain.PasswordHasher,
) *RegisterUsecase {
	return &RegisterUsecase{
		authRepo: authRepo,
		hasher:   hasher,
	}
}

type SignUpParams struct {
	Email    string
	Password string

	Name     string
	Username string
	Phone    *string
}

func (u *RegisterUsecase) Register(params SignUpParams) error {
	existing, err := u.authRepo.GetByEmail(params.Email)
	if err != nil {
		return fmt.Errorf("failed to check account: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("failed to register: account already exist")
	}

	hash, err := u.hasher.Hash(params.Password)
	if err != nil {
		return fmt.Errorf("failed to hashed: %w", err)
	}

	id := uuid.New()
	now := time.Now()

	acc := authRepo.CreateAccountProps{
		ID:           id,
		Name:         params.Name,
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: hash,
		Phone:        params.Phone,
		CreatedAt:    now,
	}

	if err := u.authRepo.Create(acc); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	return nil
}
