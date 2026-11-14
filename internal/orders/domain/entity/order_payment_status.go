package entity

type OrderPaymentStatus string

const (
	approved OrderPaymentStatus = "approved"
	rejected OrderPaymentStatus = "rejected"
	pending  OrderPaymentStatus = "pending"
)
