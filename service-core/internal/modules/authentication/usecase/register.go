package usecase

import (
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
	authDomain "service-core/internal/modules/authentication/domain"
	authRepo "service-core/internal/modules/authentication/repository"
	userRepo "service-core/internal/modules/user/repository"

	"github.com/google/uuid"
)

type RegisterUsecase struct {
	authRepo authRepo.AuthRepository
	hasher   authDomain.PasswordHasher
	userRepo userRepo.UserRepository
}

func NewRegisterUsecase(
	authRepo authRepo.AuthRepository,
	hasher authDomain.PasswordHasher,
	userRepo userRepo.UserRepository,
) *RegisterUsecase {
	return &RegisterUsecase{
		authRepo: authRepo,
		hasher:   hasher,
		userRepo: userRepo,
	}
}

type SignUpParams struct {
	Name     string
	Username string
	Email    string
	Password string
	Phone    *string
}

func (u *RegisterUsecase) Execute(params SignUpParams) error {
	now := time.Now()

	existingUser, err := u.userRepo.GetByUsername(params.Username)
	if err != nil {
		return fmt.Errorf("failed to check user: %w", err)
	}
	existingEmail, err := u.authRepo.GetByEmail(params.Email)
	if err != nil {
		return fmt.Errorf("failed to check account: %w", err)
	}

	if existingEmail != nil || existingUser != nil {
		return appErr.NewConflict(authDomain.ErrAccountAlreadyExists.Error())
	}

	user := userRepo.CreateUserProps{
		ID:        uuid.New(),
		Name:      params.Name,
		Username:  params.Username,
		Phone:     params.Phone,
		CreatedAt: now,
	}
	if err := u.userRepo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	hash, err := u.hasher.Hash(params.Password)
	if err != nil {
		return fmt.Errorf("failed to hashed: %w", err)
	}
	acc := authRepo.CreateAccountProps{
		ID:        uuid.New(),
		UserID:    user.ID,
		Email:     params.Email,
		Password:  hash,
		CreatedAt: now,
	}

	if err := u.authRepo.Create(acc); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	return nil
}
