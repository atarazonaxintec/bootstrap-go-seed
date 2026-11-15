package entity

type OrderPaymentStatus string

const (
	Approved OrderPaymentStatus = "approved"
	Rejected OrderPaymentStatus = "rejected"
	Pending  OrderPaymentStatus = "pending"
	Empty    OrderPaymentStatus = "empty"
)
