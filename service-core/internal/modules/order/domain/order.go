package domain

import (
	"fmt"
	"slices"
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

var allowedTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending: {
		OrderStatusConfirmed,
		OrderStatusCancelled,
	},

	OrderStatusConfirmed: {
		OrderStatusProcessing,
		OrderStatusCancelled,
	},

	OrderStatusProcessing: {
		OrderStatusShipped,
		OrderStatusCancelled,
	},

	OrderStatusShipped: {
		OrderStatusDelivered,
	},

	OrderStatusDelivered: {},

	OrderStatusCancelled: {},
}

type Order struct {
	ID     uuid.UUID
	Number string

	CustomerID uuid.UUID
	AddressID  uuid.UUID

	Status OrderStatus

	Subtotal    int64
	ShippingFee int64
	Total       int64

	CreatedAt time.Time
	UpdatedAt *time.Time
}

func (o *Order) UpdateStatus(status OrderStatus) error {
	if o.Status == status {
		return nil
	}

	if !o.canTransitionTo(status) {
		return fmt.Errorf("invalid status transition: %s → %s", o.Status, status)
	}

	o.Status = status
	return nil
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

func (o *Order) canTransitionTo(next OrderStatus) bool {
	allowed, exists := allowedTransitions[o.Status]
	if !exists {
		return false
	}

	return slices.Contains(allowed, next)
}

func NewOrderNumber() string {
	return fmt.Sprintf(
		"ORD-%s-%s",
		time.Now().Format("20060102"),
		strings.ToUpper(uuid.NewString()[:6]),
	)
}
