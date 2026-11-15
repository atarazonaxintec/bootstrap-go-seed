package usecase

import (
	"context"

	"microservice/internal/orders/application/dto"
	"microservice/internal/orders/domain/entity"
	domain "microservice/internal/orders/domain/repository"
)

type OrdersUseCase struct {
	repo domain.OrdersRepository
}

// Constructor del Caso de Uso
func NewOrdersUseCase(repo domain.OrdersRepository) *OrdersUseCase {
	return &OrdersUseCase{repo: repo}
}

func (uc *OrdersUseCase) NewOrders(ctx context.Context, rut string) (*dto.OrdersResponse, error) {

	responseEntity, err := uc.repo.CreateOrder(ctx, rut)
	if err != nil {
		return nil, err
	}

	response := dto.EditOrderResponse(responseEntity)

	// Aqui podriamos agregar otras validaciones:
	// - Verificar conexión a DB
	// - Verificación con servicios externos
	// Si alguno falla, se retorna error HTTP 503
	return (response), nil
}

func (uc *OrdersUseCase) UpdateOrders(ctx context.Context, order *entity.OrderEntity) (*dto.OrdersResponse, error) {
	updatedOrder, err := uc.repo.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return dto.EditOrderResponse(updatedOrder), nil
}

func (uc *OrdersUseCase) GetOrderById(ctx context.Context, id string) (*dto.OrdersResponse, error) {
	order, err := uc.repo.GetOrderById(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil || order.Id == "" {
		return nil, nil
	}

	return dto.EditOrderResponse(order), nil
}

func (uc *OrdersUseCase) GetOrdersByRut(ctx context.Context, rut string) ([]*dto.OrdersResponse, error) {
	orders, err := uc.repo.GetOrdersByRut(ctx, rut)
	if err != nil {
		return nil, err
	}

	var responses []*dto.OrdersResponse
	for _, order := range orders {
		responses = append(responses, dto.EditOrderResponse(order))
	}

	if len(responses) == 0 {
		return nil, nil
	}

	return responses, nil
}
