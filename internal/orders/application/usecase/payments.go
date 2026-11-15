package usecase

import (
	"context"
	dto "microservice/internal/orders/application/dto"
	"microservice/internal/orders/domain/entity"
	domain "microservice/internal/orders/domain/repository"
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
	payment, err := uc.paymentsRepo.CreateEmpty(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *PaymentsUseCase) FindPaymentsByOrderID(ctx context.Context, orderID string) ([]*dto.OrderPaymentResponse, error) {
	payments, err := uc.paymentsRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var responses []*dto.OrderPaymentResponse
	for _, payment := range payments {
		responses = append(responses, dto.OrderPaymentToOrderPaymentResponse(&payment))
	}

	return responses, nil
}

func (uc *PaymentsUseCase) Pay(ctx context.Context, orderID string, paymentMethod string, amount float64) (*entity.OrderPaymentEntity, error) {
	payment, err := uc.paymentsRepo.Pay(ctx, orderID, paymentMethod, amount)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (uc *PaymentsUseCase) UpdatePaymentStatus(ctx context.Context, orderID string, status entity.OrderPaymentStatus) error {
	err := uc.paymentsRepo.UpdateStatus(ctx, orderID, status)
	if err != nil {
		return err
	}

	return nil
}
