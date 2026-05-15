package usecase

import (
	"fmt"
	"time"

	appErr "service-core/internal/common/errors"
	"service-core/internal/modules/authentication/domain"
	"service-core/internal/modules/authentication/repository"

	"github.com/google/uuid"
)

type LoginEmailUsecase struct {
	authRepo repository.AuthRepository
	hasher   domain.PasswordHasher
	tokenSvc domain.TokenService
}

func NewLoginEmailUsecase(
	authRepo repository.AuthRepository,
	hasher domain.PasswordHasher,
	tokenSvc domain.TokenService,
) *LoginEmailUsecase {
	return &LoginEmailUsecase{
		authRepo: authRepo,
		hasher:   hasher,
		tokenSvc: tokenSvc,
	}
}

func (u *LoginEmailUsecase) Execute(email string, password string) (*string, time.Time, error) {
	existing, err := u.authRepo.GetByEmail(email)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to retrieve account: %w", err)
	}
	if existing == nil {
		return nil, time.Time{}, appErr.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	if err := u.hasher.Compare(existing.Password, password); err != nil {
		return nil, time.Time{}, appErr.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	token, expiry, err := u.tokenSvc.Generate(existing.ID)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to generate session: %w", err)
	}

	return &token, expiry, nil
}

func (u *LoginEmailUsecase) ById(id uuid.UUID, password string) (*string, time.Time, error) {
	existing, err := u.authRepo.GetByID(id)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to retrieve account: %w", err)
	}
	if existing == nil {
		return nil, time.Time{}, appErr.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	if err := u.hasher.Compare(existing.Password, password); err != nil {
		return nil, time.Time{}, appErr.NewUnauthorized(domain.ErrInvalidCredentials.Error())
	}

	token, expiry, err := u.tokenSvc.Generate(existing.ID)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to generate session: %w", err)
	}

	return &token, expiry, nil
}
