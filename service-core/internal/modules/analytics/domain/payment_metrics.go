package domain

import (
	"github.com/google/uuid"
)

type PaymentSummary struct {
	TotalPaid          int64
	TotalPending       int64
	TotalExpired       int64
	TotalRefunded      int64
	PaymentSuccessRate float64
	AvgTimeToPay       float64
}

type PaymentMethodBreakdown struct {
	MethodID    uuid.UUID
	MethodName  string
	MethodType  string
	Count       int
	Amount      int64
	SuccessRate float64
}
