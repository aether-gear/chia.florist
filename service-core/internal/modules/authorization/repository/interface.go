package repository

import (
	"service-core/internal/modules/authorization/domain"

	"github.com/google/uuid"
)

type RoleRepository interface {
	ListByUserID(userID uuid.UUID) ([]domain.Role, error)

	AssignToUserID(userID, roleID uuid.UUID) error

	RemoveFromUser(userID, roleID uuid.UUID)

	HasRole(userID uuid.UUID, code string) (bool, error)
}
