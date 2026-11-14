package entity

type OrderPayment struct {
	Id            string             `json:"id" database:"id"`
	OrderId       string             `json:"order_id" database:"order_id"`
	PaymentMethod string             `json:"payment_method" database:"payment_method"`
	Amount        float64            `json:"amount" database:"amount"`
	PaidAt        string             `json:"paid_at" database:"paid_at"`
	Status        OrderPaymentStatus `json:"status" database:"status"`
}
