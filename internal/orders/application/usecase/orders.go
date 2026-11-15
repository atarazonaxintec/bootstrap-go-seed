package usecase

import (
	"context"

	"microservice/internal/orders/application/dto"
	"microservice/internal/orders/domain/entity"
	domain "microservice/internal/orders/domain/repository"
	"microservice/pkg/logger"

	"go.uber.org/zap"
)

type OrdersUseCase struct {
	repo domain.OrdersRepository
}

// Constructor del Caso de Uso
func NewOrdersUseCase(repo domain.OrdersRepository) *OrdersUseCase {
	return &OrdersUseCase{repo: repo}
}

func (uc *OrdersUseCase) NewOrders(ctx context.Context, rut string) (*dto.OrdersResponse, error) {
	logger.General("Creating order")
	responseEntity, err := uc.repo.CreateOrder(ctx, rut)
	if err != nil {
		logger.Error("Failed to create order", zap.Error(err))
		return nil, err
	}

	response := dto.EditOrderResponse(responseEntity)

	// Aqui podriamos agregar otras validaciones:
	// - Verificar conexión a DB
	// - Verificación con servicios externos
	// Si alguno falla, se retorna error HTTP 503
	logger.Success("Order created successfully")
	return (response), nil
}

func (uc *OrdersUseCase) UpdateOrders(ctx context.Context, order *entity.OrderEntity) (*dto.OrdersResponse, error) {
	logger.General("Updating order")
	updatedOrder, err := uc.repo.UpdateOrder(ctx, order)
	if err != nil {
		logger.Error("Failed to update order", zap.Error(err))
		return nil, err
	}

	logger.Success("Order updated successfully")
	return dto.EditOrderResponse(updatedOrder), nil
}

func (uc *OrdersUseCase) GetOrderById(ctx context.Context, id string) (*dto.OrdersResponse, error) {
	logger.General("Fetching order by ID")
	order, err := uc.repo.GetOrderById(ctx, id)
	if err != nil {
		logger.Error("Failed to fetch order by ID", zap.Error(err))
		return nil, err
	}

	if order == nil || order.Id == "" {
		logger.General("Order not found")
		logger.Success("Order lookup completed successfully")
		return nil, nil
	}

	logger.Success("Order fetched successfully")
	return dto.EditOrderResponse(order), nil
}

func (uc *OrdersUseCase) GetOrdersByRut(ctx context.Context, rut string) ([]*dto.OrdersResponse, error) {
	logger.General("Fetching orders by RUT")
	orders, err := uc.repo.GetOrdersByRut(ctx, rut)
	if err != nil {
		logger.Error("Failed to fetch orders by RUT", zap.Error(err))
		return nil, err
	}

	var responses []*dto.OrdersResponse
	for _, order := range orders {
		responses = append(responses, dto.EditOrderResponse(order))
	}

	if len(responses) == 0 {
		logger.General("No orders found for RUT")
		logger.Success("Order lookup completed successfully")
		return nil, nil
	}

	logger.Success("Orders fetched successfully")
	return responses, nil
}
