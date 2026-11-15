package usecase

import (
	"context"

	"microservice/internal/orders/application/dto"
	domain "microservice/internal/orders/domain/repository"
	"microservice/pkg/logger"

	"go.uber.org/zap"
)

type ProductsUseCase struct {
	productsRepo domain.OrderProductsRepository
}

func NewProductsUseCase(productsRepo domain.OrderProductsRepository) *ProductsUseCase {
	return &ProductsUseCase{productsRepo: productsRepo}
}

func (uc *ProductsUseCase) FindAll(ctx context.Context) ([]*dto.ProductsResponse, error) {
	logger.General("Fetching all products")
	products, err := uc.productsRepo.FindAll(ctx)
	if err != nil {
		logger.Error("Failed to fetch products", zap.Error(err))
		return nil, err
	}

	if len(products) == 0 {
		logger.General("No products found")
		logger.Success("Product lookup completed successfully")
		return nil, nil
	}

	responses := make([]*dto.ProductsResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, dto.ProductsResponseFromEntity(product))
	}
	logger.Success("Products fetched successfully")
	return responses, nil
}

func (uc *ProductsUseCase) FindBySku(ctx context.Context, sku string) (*dto.ProductsResponse, error) {
	logger.General("Fetching product by SKU")
	product, err := uc.productsRepo.FindBySku(ctx, sku)
	if err != nil {
		logger.Error("Failed to fetch product by SKU", zap.Error(err))
		return nil, err
	}

	if product == nil {
		logger.General("Product not found")
		logger.Success("Product lookup completed successfully")
		return nil, nil
	}

	logger.Success("Product fetched successfully")
	return dto.ProductsResponseFromEntity(product), nil
}
