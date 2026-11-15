package dto

import "microservice/internal/orders/domain/entity"

type ProductsResponse struct {
	Sku           string  `json:"sku"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

func ProductsResponseFromEntity(product *entity.ProductEntity) *ProductsResponse {
	return &ProductsResponse{
		Sku:           product.Sku,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}
