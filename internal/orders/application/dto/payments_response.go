package dto

import "microservice/internal/orders/domain/entity"

type OrderPaymentResponse struct {
	PaymentID     string  `json:"payment_id" example:"123456789"`
	PaymentMethod string  `json:"payment_method" example:"credit_card"`
	Status        string  `json:"status" example:"approved"`
	Amount        float64 `json:"amount" example:"100.50"`
	PaidAt        string  `json:"paid_at" example:"2023-01-01T12:00:00Z"`
}

func OrderPaymentToOrderPaymentResponse(payment *entity.OrderPaymentEntity) *OrderPaymentResponse {
	return &OrderPaymentResponse{
		PaymentID:     payment.Id,
		PaymentMethod: payment.PaymentMethod,
		Status:        string(payment.Status),
		Amount:        payment.Amount,
		PaidAt:        payment.PaidAt,
	}
}
