package domain

import "github.com/google/uuid"

type ShopCourier struct {
	ShopID uuid.UUID
	Code   string
	Active bool
}
