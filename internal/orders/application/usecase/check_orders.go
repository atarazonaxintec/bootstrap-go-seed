package usecase

import (
	"context"
	"microservice/internal/orders/application/dto"
	domain "microservice/internal/orders/domain/repository"
)

type OrdersUseCase struct {
	repo domain.OrdersRepository
}

// Constructor del Caso de Uso
func NewOrdersUseCase(repo domain.OrdersRepository) *OrdersUseCase {
	return &OrdersUseCase{repo: repo}
}

func (uc *OrdersUseCase) Execute(ctx context.Context) (*dto.OrdersResponse, error) {

	dbStatus := "connected"

	if err := uc.repo.Ping(ctx); err != nil {
		dbStatus = "disconnected"
	}

	// Aqui podriamos agregar otras validaciones:
	// - Verificar conexión a DB
	// - Verificación con servicios externos
	// Si alguno falla, se retorna error HTTP 503
	return dto.NewOrdersResponse(dbStatus), nil
}
