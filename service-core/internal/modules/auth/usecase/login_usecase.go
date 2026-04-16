package usecase

import (
	"errors"
	"service-core/internal/modules/auth/domain"
	"service-core/internal/modules/auth/repository"
	"time"

	"github.com/google/uuid"
)

type LoginUsecase struct {
	repo     repository.AuthRepository
	hasher   domain.PasswordHasher
	tokenSvc domain.TokenService
}

func NewLoginUsecase(repo repository.AuthRepository, hasher domain.PasswordHasher, tokenSvc domain.TokenService) *LoginUsecase {
	return &LoginUsecase{
		repo:     repo,
		hasher:   hasher,
		tokenSvc: tokenSvc,
	}
}

func (u *LoginUsecase) ByEmail(email string, password string) (*string, time.Time, error) {
	existing, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, time.Time{}, err
	}

	if existing == nil {
		return nil, time.Time{}, errors.New("invalid credentials")
	}

	if err := u.hasher.Compare(existing.Password, password); err != nil {
		return nil, time.Time{}, errors.New("invalid credentials")
	}

	token, expiry, err := u.tokenSvc.Generate(existing.ID)
	if err != nil {
		return nil, time.Time{}, err
	}

	return &token, expiry, nil
}

func (u *LoginUsecase) ById(id uuid.UUID, password string) (*string, time.Time, error) {
	existing, err := u.repo.GetByID(id)
	if err != nil {
		return nil, time.Time{}, err
	}

	if existing == nil {
		return nil, time.Time{}, errors.New("invalid credentials")
	}

	if err := u.hasher.Compare(existing.Password, password); err != nil {
		return nil, time.Time{}, errors.New("invalid credentials")
	}

	token, expiry, err := u.tokenSvc.Generate(existing.ID)
	if err != nil {
		return nil, time.Time{}, err
	}

	return &token, expiry, nil
}
