package usecase

import (
	"fmt"
	authDomain "service-core/internal/features/auth/domain"
	authRepo "service-core/internal/features/auth/repository"
	"time"

	"github.com/google/uuid"
)

type SignUpParams struct {
	Email    string
	Password string

	Name     string
	Username string
	Phone    *string
}

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

func (u *RegisterUsecase) Register(params SignUpParams) error {
	existing, err := u.authRepo.GetByEmail(params.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("account already exist")
	}

	hash, err := u.hasher.Hash(params.Password)
	if err != nil {
		return err
	}

	id := uuid.New()
	now := time.Now()

	props := authRepo.CreateAccountProps{
		ID:           id,
		Name:         params.Name,
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: hash,
		Phone:        params.Phone,
		CreatedAt:    now,
	}

	if err := u.authRepo.Create(props); err != nil {
		return err
	}

	return nil
}
