package domain

import (
	"context"
	entity "microservice/internal/orders/domain/entity"
)

type OrdersRepository interface {
	Ping(ctx context.Context) error
	CreateOrder(ctx context.Context, rut string) (*entity.OrderEntity, error)
	UpdateOrder(ctx context.Context, order *entity.OrderEntity) (*entity.OrderEntity, error)
	GetOrderById(ctx context.Context, id string) (*entity.OrderEntity, error)
	GetOrdersByRut(ctx context.Context, rut string) ([]*entity.OrderEntity, error)
	GetOrderAndReserveForUpdate(ctx context.Context, id string) (*entity.OrderEntity, error)
}
