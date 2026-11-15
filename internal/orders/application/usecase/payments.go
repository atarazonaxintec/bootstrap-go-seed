package usecase

import (
	"context"

	dto "microservice/internal/orders/application/dto"
	"microservice/internal/orders/domain/entity"
	domain "microservice/internal/orders/domain/repository"
	"microservice/pkg/logger"

	"go.uber.org/zap"
)

type PaymentsUseCase struct {
	paymentsRepo domain.OrderPaymentRepository
}

func NewPaymentsUseCase(paymentsRepo domain.OrderPaymentRepository) *PaymentsUseCase {
	return &PaymentsUseCase{
		paymentsRepo: paymentsRepo,
	}
}

func (uc *PaymentsUseCase) CreateEmptyPayment(ctx context.Context, orderID string) (*entity.OrderPaymentEntity, error) {
	logger.General("Creating empty payment")
	payment, err := uc.paymentsRepo.CreateEmpty(ctx, orderID)
	if err != nil {
		logger.Error("Failed to create empty payment", zap.Error(err))
		return nil, err
	}

	logger.Success("Empty payment created successfully")
	return payment, nil
}

func (uc *PaymentsUseCase) FindPaymentsByOrderID(ctx context.Context, orderID string) ([]*dto.OrderPaymentResponse, error) {
	logger.General("Fetching payments by order ID")
	payments, err := uc.paymentsRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		logger.Error("Failed to fetch payments by order ID", zap.Error(err))
		return nil, err
	}

	var responses []*dto.OrderPaymentResponse
	for _, payment := range payments {
		responses = append(responses, dto.OrderPaymentToOrderPaymentResponse(payment))
	}

	logger.Success("Payments fetched successfully")
	return responses, nil
}

func (uc *PaymentsUseCase) Pay(ctx context.Context, orderID string, paymentMethod string, amount float64) (*entity.OrderPaymentEntity, error) {
	logger.General("Processing payment")
	payment, err := uc.paymentsRepo.Pay(ctx, orderID, paymentMethod, amount)
	if err != nil {
		logger.Error("Failed to process payment", zap.Error(err))
		return nil, err
	}

	logger.Success("Payment processed successfully")
	return payment, nil
}

func (uc *PaymentsUseCase) UpdatePaymentStatus(ctx context.Context, orderID string, status entity.OrderPaymentStatus) error {
	logger.General("Updating payment status")
	err := uc.paymentsRepo.UpdateStatus(ctx, orderID, status)
	if err != nil {
		logger.Error("Failed to update payment status", zap.Error(err))
		return err
	}

	logger.Success("Payment status updated successfully")
	return nil
}
