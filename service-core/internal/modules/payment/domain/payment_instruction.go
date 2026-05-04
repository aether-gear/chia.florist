package domain

import (
	"time"
)

type PaymentInstruction struct {
	MethodName string
	Type       string

	AccountName   string
	AccountNumber string
	PhoneNumber   string

	QRString string

	ExpiredAt time.Time
}
