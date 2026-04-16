package persistence

import "github.com/google/uuid"

type CartItemModel struct {
	ID        uuid.UUID `db:"id"`
	CartID    uuid.UUID `db:"cart_id"`
	ProductID uuid.UUID `db:"product_id"`
	Quantity  int       `db:"quantity"`
	CreatedAt int       `db:"created_at"`
	UpdatedAt int       `db:"updated_at"`
	DeletedAt int       `db:"deleted_at"`
}
