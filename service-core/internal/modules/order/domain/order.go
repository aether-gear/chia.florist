package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID     uuid.UUID
	Number string

	UserID    uuid.UUID
	AddressID uuid.UUID

	Status OrderStatus

	Subtotal    int64
	ShippingFee int64
	Total       int64

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (o Order) NewInvoice() Invoice {
	return Invoice{
		ID:          uuid.New(),
		Number:      NewInvoiceNumber(),
		OrderID:     o.ID,
		Status:      InvoiceStatusIssued,
		Subtotal:    o.Subtotal,
		ShippingFee: o.ShippingFee,
		Total:       o.Total,
		IssuedAt:    o.CreatedAt,
		CreatedAt:   o.CreatedAt,
	}
}

func NewOrderNumber() string {
	return fmt.Sprintf(
		"ORD-%s-%s",
		time.Now().Format("20060102"),
		strings.ToUpper(uuid.NewString()[:6]),
	)
}
