package repository

import "service-core/internal/modules/cart/domain"

type CartWithItems struct {
	*domain.Cart
}
