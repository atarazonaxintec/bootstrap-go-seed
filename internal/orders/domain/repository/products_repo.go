package domain

import (
	"context"
	productsEntity "microservice/internal/orders/domain/entity"
)

type OrderProductsRepository interface {
	FindAll(ctx context.Context) ([]*productsEntity.ProductEntity, error)
	FindBySku(ctx context.Context, sku string) (*productsEntity.ProductEntity, error)
}
