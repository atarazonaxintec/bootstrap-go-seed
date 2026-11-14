package domain

import (
	"context"
	paymentsEntity "microservice/internal/orders/domain/entity"
	statuses "microservice/internal/orders/domain/entity"
)

type OrderPaymentRepository interface {
	CreateEmpty(ctx context.Context, orderID string) (*paymentsEntity.OrderPayment, error)
	FindByOrderID(ctx context.Context, orderID string) ([]paymentsEntity.OrderPayment, error)
	Pay(ctx context.Context, orderID string, paymentMethod string, amount float64) (*paymentsEntity.OrderPayment, error)
	UpdateStatus(ctx context.Context, orderID string, status statuses.OrderPaymentStatus) error
}
