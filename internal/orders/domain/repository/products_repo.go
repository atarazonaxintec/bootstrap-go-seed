package domain

import (
	"context"
	productsEntity "microservice/internal/orders/domain/entity"
)

type ProductsRepository interface {
	FindAll(ctx context.Context) ([]productsEntity.Product, error)
	FindBySku(ctx context.Context, sku string) (*productsEntity.Product, error)
}
