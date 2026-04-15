package persistence

import (
	"service-core/internal/modules/auth/domain"
)

func (m *AccountModel) ToDomain() (*domain.Account, error) {
	return &domain.Account{
		ID:          m.ID,
		Email:       m.Email,
		Password:    m.Password,
		LastLoginAt: m.LastLoginAt,
	}, nil
}

func FromDomain(p *domain.Account) *AccountModel {
	return &AccountModel{
		ID:          p.ID,
		Email:       p.Email,
		Password:    p.Password,
		LastLoginAt: p.LastLoginAt,
	}
}
