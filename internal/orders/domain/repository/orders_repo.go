package domain

import (
	"context"
	entity "microservice/internal/orders/domain/entity"
)

type OrdersRepository interface {
	Ping(ctx context.Context) error
	CreateOrder(ctx context.Context, rut string) (*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetOrderById(ctx context.Context, id string) (*entity.Order, error)
	GetOrdersByRut(ctx context.Context, rut string) ([]*entity.Order, error)
}
