package usecase

import (
	"context"

	"microservice/internal/orders/application/dto"
	domain "microservice/internal/orders/domain/repository"
)

type ProductsUseCase struct {
	productsRepo domain.OrderProductsRepository
}

func NewProductsUseCase(productsRepo domain.OrderProductsRepository) *ProductsUseCase {
	return &ProductsUseCase{productsRepo: productsRepo}
}

func (uc *ProductsUseCase) FindAll(ctx context.Context) ([]*dto.ProductsResponse, error) {
	products, err := uc.productsRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	responses := make([]*dto.ProductsResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, dto.ProductsResponseFromEntity(product))
	}
	return responses, nil
}

func (uc *ProductsUseCase) FindBySku(ctx context.Context, sku string) (*dto.ProductsResponse, error) {
	product, err := uc.productsRepo.FindBySku(ctx, sku)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, nil
	}

	return dto.ProductsResponseFromEntity(product), nil
}
