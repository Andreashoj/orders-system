package models

import "time"

type PaymentMethod string

const (
	Paypal    PaymentMethod = "PayPal"
	Visa      PaymentMethod = "Visa"
	MobilePay PaymentMethod = "MobilePay"
)

type Transaction struct {
	ID          int
	OrderID     int
	PaymentType PaymentMethod
	CreatedAt   time.Time
}
