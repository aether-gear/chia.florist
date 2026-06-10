package domain

import "github.com/google/uuid"

type Permission struct {
	Id   uuid.UUID
	Code string
}
